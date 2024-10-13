package article

import (
	"github.com/core-go/core"
	s "github.com/core-go/search/handler"
)

func NewArticleHandler(search s.Search[Article, *ArticleFilter], service ArticleService, logError core.Log, validate core.Validate[*Article], writeLog core.WriteLog, action *core.ActionConfig) ArticleTransport {
	hdl := core.NewhandlerWithLog[Article, string](service, logError, validate, action, writeLog)
	searchHandler := s.NewSearchHandler[Article, *ArticleFilter](search, logError, nil)
	return &ArticleHandler{Handler: hdl, SearchHandler: searchHandler}
}

type ArticleHandler struct {
	*core.Handler[Article, string]
	*s.SearchHandler[Article, *ArticleFilter]
}
