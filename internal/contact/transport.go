package contact

import (
	"database/sql"
	"net/http"

	"github.com/core-go/core"
	v "github.com/core-go/core/validator"
	"github.com/core-go/sql/query/builder"
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
	contactRepository, err := NewContactAdapter(db, queryContact)
	if err != nil {
		return nil, err
	}
	contactService := NewContactService(db, contactRepository)
	contactHandler := NewContactHandler(contactService, logError, validator.Validate, action)
	return contactHandler, nil
}
