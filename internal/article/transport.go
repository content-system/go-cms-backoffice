package article

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/core-go/core"
	"github.com/core-go/core/approver"
	notification "github.com/core-go/core/notification/adapter"
	v "github.com/core-go/core/validator"
	q "github.com/core-go/sql"
	"github.com/lib/pq"
	"github.com/teris-io/shortid"

	histories "github.com/core-go/core/histories/adapter"
	history "github.com/core-go/core/history/adapter"
)

type ArticleTransport interface {
	Search(w http.ResponseWriter, r *http.Request)
	LoadDraft(w http.ResponseWriter, r *http.Request)
	Load(w http.ResponseWriter, r *http.Request)
	GetHistories(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Approve(w http.ResponseWriter, r *http.Request)
	Reject(w http.ResponseWriter, r *http.Request)
	Patch(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

func NewArticleTransport(db *sql.DB, logError core.Log, writeLog core.WriteLog, action *core.ActionConfig) (ArticleTransport, error) {
	validator, err := v.NewValidator[*Article]()
	if err != nil {
		return nil, err
	}

	draftArticleRepository, err := NewDraftArticleAdapter(db, BuildQuery, pq.Array)
	if err != nil {
		return nil, err
	}

	articleRepository, err := NewArticleAdapter(db, pq.Array)
	if err != nil {
		return nil, err
	}

	buildParam := q.GetBuild(db)
	approverPort := approver.NewApproversAdapter(db, "article", sqlGetApprover)
	notificationPort := notification.NewNotificationAdapter(db, Generate, buildParam, "notifications", "tx", "time")
	historyPort := history.NewHistoryAdapter(db, Generate, buildParam, "entity", "article", "histories", "author", "time")
	articleService := NewArticleService(db, draftArticleRepository, articleRepository, historyPort, approverPort, notificationPort)

	historyQuery := histories.NewHistoryAdapter(db, buildParam, nil, "histories", "entity", "author", "time")
	articleHandler := NewArticleHandler(articleService, logError, validator.Validate, writeLog, action, historyQuery)

	return articleHandler, nil
}

var sid *shortid.Shortid

func Generate(ctx context.Context) (string, error) {
	if sid == nil {
		s, err := shortid.New(1, shortid.DefaultABC, 2342)
		if err != nil {
			return "", err
		}

		sid = s
	}

	return sid.Generate()
}

const sqlGetApprover = `
	select
		distinct u.user_id
	from
		users u
	join user_roles ur on
		u.user_id = ur.user_id
	join role_modules r on
		ur.role_id = r.role_id
	where
		(r.permissions & 8) = 8
		and u.status = 'A'
		and r.module_id = $1
`
