package article

import (
	"database/sql"
	"net/http"

	"github.com/core-go/core"
	sv "github.com/core-go/core/sql"
	v "github.com/core-go/core/validator"
	"github.com/core-go/sql/adapter"
	"github.com/core-go/sql/query/builder"
	"github.com/lib/pq"
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
	queryArticle := builder.UseQuery[Article, *ArticleFilter](db, "news")
	articleSearchBuilder, err := adapter.NewSearchAdapterWithArray[Article, string, *ArticleFilter](db, "news", queryArticle, pq.Array)
	if err != nil {
		return nil, err
	}
	articleRepository, err := adapter.NewAdapterWithArray[Article, string](db, "news", pq.Array)
	if err != nil {
		return nil, err
	}
	articleService := sv.NewService[Article, string](db, articleRepository)
	articleHandler := NewArticleHandler(articleSearchBuilder.Search, articleService, logError, validator.Validate, writeLog, action)
	return articleHandler, nil
}
