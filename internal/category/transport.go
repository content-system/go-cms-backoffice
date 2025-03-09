package category

import (
	"database/sql"
	"net/http"

	"github.com/core-go/core"
	v "github.com/core-go/core/validator"
	"github.com/core-go/sql/query/builder"
)

type CategoryTransport interface {
	Search(w http.ResponseWriter, r *http.Request)
	Load(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Patch(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

func NewCategoryTransport(db *sql.DB, logError core.Log, writeLog core.WriteLog, action *core.ActionConfig) (CategoryTransport, error) {
	validator, err := v.NewValidator[*Category]()
	if err != nil {
		return nil, err
	}
	queryCategory := builder.UseQuery[Category, *CategoryFilter](db, "categories")
	categoryRepository, err := NewCategoryAdapter(db, queryCategory)
	if err != nil {
		return nil, err
	}
	categoryService := NewCategoryService(db, categoryRepository)
	categoryHandler := NewCategoryHandler(categoryService, logError, validator.Validate, writeLog, action)
	return categoryHandler, nil
}
