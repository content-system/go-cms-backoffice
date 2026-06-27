package article

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/core-go/core"
	"github.com/core-go/core/approver"
	"github.com/core-go/core/history"
	"github.com/core-go/core/notification"
	"github.com/core-go/core/shortid"
	"github.com/core-go/core/tx"

	act "github.com/core-go/core/action"
	"github.com/core-go/core/slug"
	"github.com/core-go/core/status"
)

type ArticleService interface {
	LoadDraft(ctx context.Context, id string) (*Article, error)
	Load(ctx context.Context, id string) (*Article, error)
	Create(ctx context.Context, article *Article) (int64, error)
	Update(ctx context.Context, article *Article) (int64, error)
	Patch(ctx context.Context, article map[string]interface{}) (int64, error)
	Delete(ctx context.Context, id string) (int64, error)
	Search(ctx context.Context, filter *ArticleFilter, limit int64, offset int64) ([]Article, int64, error)
	Approve(ctx context.Context, id string, approvedBy string) (int64, error)
	Reject(ctx context.Context, id string, rejectedBy string) (int64, error)
}

func NewArticleService(db *sql.DB, draftRepository DraftArticleRepository, repository ArticleRepository, historyRepository history.HistoryPort, approverPort approver.ApproversPort, notificationPort notification.NotificationPort) *ArticleUseCase {
	return &ArticleUseCase{db: db, draftRepository: draftRepository, repository: repository, historyRepository: historyRepository, approverPort: approverPort, notificationPort: notificationPort}
}

type ArticleUseCase struct {
	db                *sql.DB
	draftRepository   DraftArticleRepository
	repository        ArticleRepository
	historyRepository history.HistoryPort
	approverPort      approver.ApproversPort
	notificationPort  notification.NotificationPort
}

func canReject(s string) bool {
	return s == status.Submitted || s == status.Approved
}

func canUpdate(s string) bool {
	return s != status.Approved && s != status.Expired
}

func (s *ArticleUseCase) Search(ctx context.Context, filter *ArticleFilter, limit int64, offset int64) ([]Article, int64, error) {
	return s.draftRepository.Search(ctx, filter, limit, offset)
}
func (s *ArticleUseCase) LoadDraft(ctx context.Context, id string) (*Article, error) {
	return s.draftRepository.Load(ctx, id)
}
func (s *ArticleUseCase) Load(ctx context.Context, id string) (*Article, error) {
	return s.repository.Load(ctx, id)
}

func (s *ArticleUseCase) Create(ctx context.Context, article *Article) (int64, error) {
	id, err := shortid.Generate(ctx)
	if err != nil {
		return -1, err
	}

	article.Id = id
	article.Slug = slug.Slugify(article.Title, article.Id, 10, 60)
	article.AuthorId = article.CreatedBy

	action := act.Create

	if article.Status == status.Submitted {
		article.SubmittedBy = article.CreatedBy
		article.SubmittedAt = core.Now()
		action = act.Submit
	}

	res, err := tx.Execute(ctx, s.db, func(ctx context.Context) (int64, error) {

		res, err := s.draftRepository.Create(ctx, article)
		if err != nil {
			return 0, err
		}

		_, err = s.historyRepository.Create(ctx, id, article.CreatedBy, action, article.GetData())

		if err != nil {
			fmt.Println(err)
			return 0, err
		}

		if article.Status == status.Submitted {
			s.notifyApprovers(ctx, article.Id, article.SubmittedBy)
		}

		return res, nil
	})

	return res, err
}
func (s *ArticleUseCase) Update(ctx context.Context, article *Article) (int64, error) {
	return tx.Execute(ctx, s.db, func(ctx context.Context) (int64, error) {

		// 1. Check exist in main table
		isExist, err := s.repository.Exist(ctx, article.Id)
		if err != nil {
			return 0, err
		}

		// If not exists -> generate slug
		if !isExist {
			article.Slug = slug.Slugify(article.Title, article.Id, 10, 60)
		}

		// 2. Load existing draft
		existingArticle, err := s.draftRepository.Load(ctx, article.Id)
		if err != nil {
			return 0, err
		}

		if existingArticle == nil {
			return 0, nil
		}

		// 3. Validate canUpdate
		if !canUpdate(existingArticle.Status) {
			return -1, nil
		}

		// 4. Handle submit flow
		if article.Status == status.Submitted {
			article.SubmittedBy = article.UpdatedBy
			article.SubmittedAt = core.Now()
		}

		// 5. Update draft
		res, err := s.draftRepository.Update(ctx, article)
		if err != nil {
			return 0, err
		}

		// 6. if status is Submitted -> create history + notify
		if article.Status == status.Submitted {
			_, err = s.historyRepository.Create(ctx, article.Id, article.UpdatedBy, act.Submit, article.GetData())
			if err != nil {
				return 0, err
			}
			s.notifyApprovers(ctx, article.Id, article.SubmittedBy)
		}

		return res, nil
	})
}
func (s *ArticleUseCase) notifyApprovers(ctx context.Context, id string, userId string) error {
	approvers, err := s.approverPort.GetApprovers(ctx)
	if len(approvers) == 0 || err != nil {
		return err
	}

	var notifications []notification.Notification

	url := fmt.Sprintf("/articles/%s/approve", id)
	msg := fmt.Sprintf("Please review and approve an article (id: '%s').", id)
	for _, approver := range approvers {
		noti := notification.Build(userId, approver, url, msg)
		notifications = append(notifications, *noti)
	}
	_, err = s.notificationPort.PushNotifications(ctx, notifications)

	return err
}

func (s *ArticleUseCase) Approve(ctx context.Context, id string, approvedBy string) (int64, error) {

	return tx.Execute(ctx, s.db, func(ctx context.Context) (int64, error) {

		// 1. Load draft
		article, err := s.draftRepository.Load(ctx, id)
		if err != nil {
			return 0, err
		}

		if article == nil {
			return 0, nil
		}

		// 2. Validate status
		if article.Status != status.Submitted {
			return -1, nil
		}

		// 3. Maker-checker (not allow approving own submission)
		if article.SubmittedBy == approvedBy {
			return -2, nil
		}

		// 4. Update fields
		article.Status = status.Published
		article.ApprovedBy = approvedBy
		article.ApprovedAt = core.Now()

		article.UpdatedBy = approvedBy
		article.UpdatedAt = article.ApprovedAt
		article.PublishedAt = article.ApprovedAt

		// 5. Update draft
		_, err = s.draftRepository.Update(ctx, article)
		if err != nil {
			return 0, err
		}

		// 6. Save to main repository
		res, err := s.repository.Save(ctx, article)
		if err != nil {
			return 0, err
		}

		// 7. History
		_, err = s.historyRepository.Create(ctx, id, approvedBy, act.Approve, article.GetData())
		if err != nil {
			return 0, err
		}

		// 8. Notify submitter
		msg := fmt.Sprintf("This article was approved (id: '%s').", id)

		s.notifySubmitter(ctx, id, approvedBy, article.SubmittedBy, msg)

		return res, nil
	})
}
func (s *ArticleUseCase) Reject(ctx context.Context, id string, rejectedBy string) (int64, error) {
	return tx.Execute(ctx, s.db, func(ctx context.Context) (int64, error) {

		// 1. Load draft
		article, err := s.draftRepository.Load(ctx, id)
		if err != nil {
			return 0, err
		}

		if article == nil {
			return 0, nil
		}

		// 2. Validate canReject
		if !canReject(article.Status) {
			return -1, nil
		}

		// 3. Maker-checker
		if article.SubmittedBy == rejectedBy {
			return -2, nil
		}

		// 4. Update fields
		article.Status = status.Rejected
		article.ApprovedBy = rejectedBy
		article.ApprovedAt = core.Now()

		article.UpdatedBy = rejectedBy
		article.UpdatedAt = article.ApprovedAt

		// 5. Update draft
		res, err := s.draftRepository.Update(ctx, article)
		if err != nil {
			return 0, err
		}

		// 6. History
		_, err = s.historyRepository.Create(ctx, id, rejectedBy, act.Reject, article.GetData())
		if err != nil {
			return 0, err
		}

		// 7. Notify
		msg := fmt.Sprintf("This article was rejected (id: '%s').", id)

		s.notifySubmitter(ctx, id, rejectedBy, article.SubmittedBy, msg)

		return res, nil
	})
}
func (s *ArticleUseCase) notifySubmitter(ctx context.Context, id string, userId string, submitter string, message string) error {
	noti := notification.Build(userId, submitter, fmt.Sprintf("/articles/%s", id), message)
	_, err := s.notificationPort.Push(ctx, noti)
	return err
}

func (s *ArticleUseCase) Delete(ctx context.Context, id string) (int64, error) {
	return tx.Execute(ctx, s.db, func(ctx context.Context) (int64, error) {
		return s.draftRepository.Delete(ctx, id)
	})
}

func (s *ArticleUseCase) Patch(ctx context.Context, article map[string]interface{}) (int64, error) {
	return tx.Execute(ctx, s.db, func(ctx context.Context) (int64, error) {
		return s.draftRepository.Patch(ctx, article)
	})
}
