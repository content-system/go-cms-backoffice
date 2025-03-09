package category

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/core-go/core"
	"github.com/core-go/search"
)

func NewCategoryHandler(service CategoryService, logError core.Log, validate core.Validate[*Category], writeLog core.WriteLog, action *core.ActionConfig) *CategoryHandler {
	categoryType := reflect.TypeOf(Category{})
	parameters := search.CreateParameters(reflect.TypeOf(CategoryFilter{}), categoryType)
	attributes := core.CreateAttributes(categoryType, logError, writeLog, action)
	return &CategoryHandler{service: service, Validate: validate, Attributes: attributes, Parameters: parameters}
}

type CategoryHandler struct {
	service  CategoryService
	Validate core.Validate[*Category]
	*core.Attributes
	*search.Parameters
}

func (h *CategoryHandler) Load(w http.ResponseWriter, r *http.Request) {
	id, err := core.GetRequiredString(w, r)
	if err == nil {
		category, err := h.service.Load(r.Context(), id)
		if err != nil {
			h.Error(r.Context(), fmt.Sprintf("Error to get category '%s': %s", id, err.Error()))
			http.Error(w, core.InternalServerError, http.StatusInternalServerError)
			return
		}
		core.JSON(w, core.IsFound(category), category)
	}
}
func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	category, er1 := core.Decode[Category](w, r)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &category)
		if !core.HasError(w, r, errors, er2, h.Error, &category, h.Log, h.Resource, h.Action.Create) {
			res, er3 := h.service.Create(r.Context(), &category)
			if er3 != nil {
				h.Error(r.Context(), er3.Error())
				h.Log(r.Context(), h.Resource, h.Action.Update, false, er3.Error())
				http.Error(w, core.InternalServerError, http.StatusInternalServerError)
				return
			}

			if res > 0 {
				h.Log(r.Context(), h.Resource, h.Action.Update, true, fmt.Sprintf("create '%s'", category.Id))
				core.JSON(w, http.StatusCreated, category)
			} else {
				h.Log(r.Context(), h.Resource, h.Action.Update, false, fmt.Sprintf("conflict '%s'", category.Id))
				core.JSON(w, http.StatusConflict, res)
			}
		}
	}
}
func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	category, er1 := core.DecodeAndCheckId[Category](w, r, h.Keys, h.Indexes)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &category)
		if !core.HasError(w, r, errors, er2, h.Error, &category, h.Log, h.Resource, h.Action.Update) {
			res, err := h.service.Update(r.Context(), &category)
			if err != nil {
				h.Error(r.Context(), err.Error())
				h.Log(r.Context(), h.Resource, h.Action.Update, false, err.Error())
				http.Error(w, core.InternalServerError, http.StatusInternalServerError)
				return
			}

			if res > 0 {
				h.Log(r.Context(), h.Resource, h.Action.Update, true, fmt.Sprintf("%s '%s'", h.Action.Update, category.Id))
				core.JSON(w, http.StatusOK, category)
			} else if res == 0 {
				h.Log(r.Context(), h.Resource, h.Action.Update, false, fmt.Sprintf("not found '%s'", category.Id))
				core.JSON(w, http.StatusNotFound, res)
			} else {
				h.Log(r.Context(), h.Resource, h.Action.Update, false, fmt.Sprintf("conflict '%s'", category.Id))
				core.JSON(w, http.StatusConflict, res)
			}
		}
	}
}
func (h *CategoryHandler) Patch(w http.ResponseWriter, r *http.Request) {
	r, category, jsonCategory, er1 := core.BuildMapAndCheckId[Category](w, r, h.Keys, h.Indexes)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &category)
		if !core.HasError(w, r, errors, er2, h.Error, jsonCategory, h.Log, h.Resource, h.Action.Patch) {
			res, err := h.service.Patch(r.Context(), jsonCategory)
			if err != nil {
				h.Error(r.Context(), err.Error())
				h.Log(r.Context(), h.Resource, h.Action.Patch, false, err.Error())
				http.Error(w, core.InternalServerError, http.StatusInternalServerError)
				return
			}

			if res > 0 {
				h.Log(r.Context(), h.Resource, h.Action.Patch, true, fmt.Sprintf("%s '%s'", h.Action.Patch, category.Id))
				core.JSON(w, http.StatusOK, jsonCategory)
			} else if res == 0 {
				h.Log(r.Context(), h.Resource, h.Action.Patch, false, fmt.Sprintf("not found '%s'", category.Id))
				core.JSON(w, http.StatusNotFound, res)
			} else {
				h.Log(r.Context(), h.Resource, h.Action.Patch, false, fmt.Sprintf("conflict '%s'", category.Id))
				core.JSON(w, http.StatusConflict, res)
			}
		}
	}
}
func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := core.GetRequiredString(w, r)
	if err == nil {
		res, err := h.service.Delete(r.Context(), id)
		if err != nil {
			h.Error(r.Context(), err.Error())
			h.Log(r.Context(), h.Resource, h.Action.Delete, false, err.Error())
			http.Error(w, core.InternalServerError, http.StatusInternalServerError)
			return
		}

		if res > 0 {
			h.Log(r.Context(), h.Resource, h.Action.Delete, true, fmt.Sprintf("%s '%s'", h.Action.Delete, id))
			core.JSON(w, http.StatusOK, res)
		} else if res == 0 {
			h.Log(r.Context(), h.Resource, h.Action.Delete, false, fmt.Sprintf("not found '%s'", id))
			core.JSON(w, http.StatusNotFound, res)
		} else {
			h.Log(r.Context(), h.Resource, h.Action.Delete, false, fmt.Sprintf("conflict '%s'", id))
			core.JSON(w, http.StatusConflict, res)
		}
	}
}
func (h *CategoryHandler) Search(w http.ResponseWriter, r *http.Request) {
	filter := CategoryFilter{Filter: &search.Filter{}}
	err := search.Decode(r, &filter, h.ParamIndex, h.FilterIndex)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	offset := search.GetOffset(filter.Limit, filter.Page)
	categories, total, err := h.service.Search(r.Context(), &filter, filter.Limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	core.JSON(w, http.StatusOK, &search.Result{List: &categories, Total: total})
}
