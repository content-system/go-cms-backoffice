package content

import (
	"database/sql"
	"github.com/lib/pq"
	"net/http"

	"github.com/core-go/core"
	v "github.com/core-go/core/validator"
	"github.com/core-go/sql/query/builder"
)

type ContentTransport interface {
	Search(w http.ResponseWriter, r *http.Request)
	Load(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Patch(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

func NewContentTransport(db *sql.DB, logError core.Log, writeLog core.WriteLog, action *core.ActionConfig) (ContentTransport, error) {
	validator, err := v.NewValidator[*Content]()
	if err != nil {
		return nil, err
	}
	queryContent := builder.UseQuery[Content, *ContentFilter](db, "contents")
	contentRepository, err := NewContentAdapter(db, queryContent, pq.Array)
	if err != nil {
		return nil, err
	}
	contentService := NewContentService(db, contentRepository)
	contentHandler := NewContentHandler(contentService, logError, validator.Validate, writeLog, action)
	return contentHandler, nil
}
