package article

import (
	"context"
	"database/sql"
	"time"

	"github.com/core-go/core/shortid"
	"github.com/core-go/core/tx"

	"go-service/pkg/history"
	"go-service/pkg/slug"
	"go-service/pkg/status"
)

type ArticleService interface {
	LoadDraft(ctx context.Context, id string) (*Article, error)
	Load(ctx context.Context, id string) (*Article, error)
	Create(ctx context.Context, article *Article) (int64, error)
	Update(ctx context.Context, article *Article) (int64, error)
	Patch(ctx context.Context, article map[string]interface{}) (int64, error)
	Delete(ctx context.Context, id string) (int64, error)
	Search(ctx context.Context, filter *ArticleFilter, limit int64, offset int64) ([]Article, int64, error)
}

func NewArticleService(db *sql.DB, draftRepository DraftArticleRepository, repository ArticleRepository, historyRepository history.HistoryPort) *ArticleUseCase {
	return &ArticleUseCase{db: db, draftRepository: draftRepository, repository: repository, historyRepository: historyRepository}
}

type ArticleUseCase struct {
	db                *sql.DB
	repository        ArticleRepository
	draftRepository   DraftArticleRepository
	historyRepository history.HistoryPort
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
	now := time.Now()
	id, err := shortid.Generate(ctx)
	if err != nil {
		return -1, err
	}
	article.Id = id
	article.Slug = slug.Slugify(article.Title, article.Id, 10, 60)
	article.AuthorId = article.CreatedBy

	if article.Status == status.Submitted {
		article.SubmittedBy = article.UpdatedBy
		article.SubmittedAt = &now
	}

	res, err := tx.Execute(ctx, s.db, func(ctx context.Context) (int64, error) {
		return s.repository.Create(ctx, article)
	})
	if res > 0 && err == nil && article.Status == status.Submitted {

	}
	return res, err
}
func (s *ArticleUseCase) Update(ctx context.Context, article *Article) (int64, error) {
	return tx.Execute(ctx, s.db, func(ctx context.Context) (int64, error) {
		return s.repository.Update(ctx, article)
	})
}
func (s *ArticleUseCase) Patch(ctx context.Context, article map[string]interface{}) (int64, error) {
	return tx.Execute(ctx, s.db, func(ctx context.Context) (int64, error) {
		return s.repository.Patch(ctx, article)
	})
}
func (s *ArticleUseCase) Delete(ctx context.Context, id string) (int64, error) {
	return tx.Execute(ctx, s.db, func(ctx context.Context) (int64, error) {
		return s.repository.Delete(ctx, id)
	})
}
