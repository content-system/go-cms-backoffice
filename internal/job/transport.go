package job

import (
	"database/sql"
	"github.com/core-go/core"
	sv "github.com/core-go/core/sql"
	v "github.com/core-go/core/validator"
	"github.com/core-go/sql/adapter"
	"github.com/core-go/sql/query/builder"
	"net/http"
)

type JobTransport interface {
	Search(w http.ResponseWriter, r *http.Request)
	Load(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Patch(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

func NewJobTransport(db *sql.DB, logError core.Log, writeLog core.WriteLog, action *core.ActionConfig) (JobTransport, error) {
	validator, err := v.NewValidator[*Job]()
	if err != nil {
		return nil, err
	}
	queryJob := builder.UseQuery[Job, *JobFilter](db, "jobs")
	jobSearchBuilder, err := adapter.NewSearchAdapter[Job, string, *JobFilter](db, "jobs", queryJob)
	if err != nil {
		return nil, err
	}
	jobRepository, err := adapter.NewAdapter[Job, string](db, "jobs")
	if err != nil {
		return nil, err
	}
	jobService := sv.NewService[Job, string](db, jobRepository)
	jobHandler := NewJobHandler(jobSearchBuilder.Search, jobService, logError, validator.Validate, writeLog, action)
	return jobHandler, nil
}
