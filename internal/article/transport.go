package article

import (
	"database/sql"
	"github.com/lib/pq"
	"net/http"

	"github.com/core-go/core"
	v "github.com/core-go/core/validator"
	"github.com/core-go/sql/query/builder"
)

type ArticleTransport interface {
	Search(w http.ResponseWriter, r *http.Request)
	Load(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Patch(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

func NewArticleTransport(db *sql.DB, logError core.Log, writeLog core.WriteLog, action *core.ActionConfig) (ArticleTransport, error) {
	validator, err := v.NewValidator[*Article]()
	if err != nil {
		return nil, err
	}
	queryArticle := builder.UseQuery[Article, *ArticleFilter](db, "articles")
	articleRepository, err := NewArticleAdapter(db, queryArticle, pq.Array)
	if err != nil {
		return nil, err
	}
	articleService := NewArticleService(db, articleRepository)
	articleHandler := NewArticleHandler(articleService, logError, validator.Validate, writeLog, action)
	return articleHandler, nil
}
