package job

import (
	"github.com/core-go/core"
	s "github.com/core-go/search/handler"
)

func NewJobHandler(search s.Search[Job, *JobFilter], service JobService, logError core.Log, validate core.Validate[*Job], writeLog core.WriteLog, action *core.ActionConfig) JobTransport {
	hdl := core.NewhandlerWithLog[Job, string](service, logError, validate, action, writeLog)
	searchHandler := s.NewSearchHandler[Job, *JobFilter](search, logError, nil)
	return &JobHandler{Handler: hdl, SearchHandler: searchHandler}
}

type JobHandler struct {
	*core.Handler[Job, string]
	*s.SearchHandler[Job, *JobFilter]
}
