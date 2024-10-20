package contact

import (
	"github.com/core-go/core"
	s "github.com/core-go/search/handler"
)

func NewContactHandler(search s.Search[Contact, *ContactFilter], service ContactService, logError core.Log, validate core.Validate[*Contact], writeLog core.WriteLog, action *core.ActionConfig) ContactTransport {
	hdl := core.NewhandlerWithLog[Contact, string](service, logError, validate, action, writeLog)
	searchHandler := s.NewSearchHandler[Contact, *ContactFilter](search, logError, nil)
	return &ContactHandler{Handler: hdl, SearchHandler: searchHandler}
}

type ContactHandler struct {
	*core.Handler[Contact, string]
	*s.SearchHandler[Contact, *ContactFilter]
}
