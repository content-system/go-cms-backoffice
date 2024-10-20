package contact

import (
	"database/sql"
	"github.com/core-go/core"
	sv "github.com/core-go/core/sql"
	v "github.com/core-go/core/validator"
	"github.com/core-go/sql/adapter"
	"github.com/core-go/sql/query/builder"
	"net/http"
)

type ContactTransport interface {
	Search(w http.ResponseWriter, r *http.Request)
	Load(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Patch(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

func NewContactTransport(db *sql.DB, logError core.Log, writeLog core.WriteLog, action *core.ActionConfig) (ContactTransport, error) {
	validator, err := v.NewValidator[*Contact]()
	if err != nil {
		return nil, err
	}
	queryContact := builder.UseQuery[Contact, *ContactFilter](db, "contacts")
	contactSearchBuilder, err := adapter.NewSearchAdapter[Contact, string, *ContactFilter](db, "contacts", queryContact)
	if err != nil {
		return nil, err
	}
	contactRepository, err := adapter.NewAdapter[Contact, string](db, "contacts")
	if err != nil {
		return nil, err
	}
	contactService := sv.NewService[Contact, string](db, contactRepository)
	contactHandler := NewContactHandler(contactSearchBuilder.Search, contactService, logError, validator.Validate, writeLog, action)
	return contactHandler, nil
}
