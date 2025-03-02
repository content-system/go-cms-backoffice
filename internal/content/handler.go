package content

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/core-go/core"
	"github.com/core-go/search"
)

func NewContentHandler(service ContentService, logError core.Log, validate core.Validate[*Content], writeLog core.WriteLog, action *core.ActionConfig) *ContentHandler {
	contentType := reflect.TypeOf(Content{})
	parameters := search.CreateParameters(reflect.TypeOf(ContentFilter{}), contentType)
	attributes := core.CreateAttributes(contentType, logError, action, writeLog)
	return &ContentHandler{service: service, Validate: validate, Attributes: attributes, Parameters: parameters}
}

type ContentHandler struct {
	service  ContentService
	Validate core.Validate[*Content]
	*core.Attributes
	*search.Parameters
}

func (h *ContentHandler) Load(w http.ResponseWriter, r *http.Request) {
	id, er1 := core.GetRequiredString(w, r, 1)
	lang, er2 := core.GetRequiredString(w, r)
	if er1 == nil && er2 == nil {
		content, err := h.service.Load(r.Context(), id, lang)
		if err != nil {
			h.Error(r.Context(), fmt.Sprintf("Error to get content '%s': %s", id, err.Error()))
			http.Error(w, core.InternalServerError, http.StatusInternalServerError)
			return
		}
		core.JSON(w, core.IsFound(content), content)
	}
}
func (h *ContentHandler) Create(w http.ResponseWriter, r *http.Request) {
	content, er1 := core.Decode[Content](w, r)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &content)
		if !core.HasError(w, r, errors, er2, h.Error, &content, h.Log, h.Resource, h.Action.Create) {
			res, er3 := h.service.Create(r.Context(), &content)
			if er3 != nil {
				h.Error(r.Context(), er3.Error())
				h.Log(r.Context(), h.Resource, h.Action.Update, false, er3.Error())
				http.Error(w, core.InternalServerError, http.StatusInternalServerError)
				return
			}

			if res > 0 {
				h.Log(r.Context(), h.Resource, h.Action.Create, true, fmt.Sprintf("%s '%s' '%s'", h.Action.Create, content.Id, content.Lang))
				core.JSON(w, http.StatusCreated, content)
			} else {
				h.Log(r.Context(), h.Resource, h.Action.Create, false, fmt.Sprintf("conflict '%s' '%s'", content.Id, content.Lang))
				core.JSON(w, http.StatusConflict, res)
			}
		}
	}
}
func (h *ContentHandler) Update(w http.ResponseWriter, r *http.Request) {
	content, er1 := core.DecodeAndCheckId[Content](w, r, h.Keys, h.Indexes)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &content)
		if !core.HasError(w, r, errors, er2, h.Error, &content, h.Log, h.Resource, h.Action.Update) {
			res, err := h.service.Update(r.Context(), &content)
			if err != nil {
				h.Error(r.Context(), err.Error())
				h.Log(r.Context(), h.Resource, h.Action.Update, false, err.Error())
				http.Error(w, core.InternalServerError, http.StatusInternalServerError)
				return
			}

			if res > 0 {
				h.Log(r.Context(), h.Resource, h.Action.Update, true, fmt.Sprintf("%s '%s' '%s'", h.Action.Update, content.Id, content.Lang))
				core.JSON(w, http.StatusOK, content)
			} else if res == 0 {
				h.Log(r.Context(), h.Resource, h.Action.Update, false, fmt.Sprintf("not found '%s' '%s'", content.Id, content.Lang))
				core.JSON(w, http.StatusNotFound, res)
			} else {
				h.Log(r.Context(), h.Resource, h.Action.Update, false, fmt.Sprintf("conflict '%s' '%s'", content.Id, content.Lang))
				core.JSON(w, http.StatusConflict, res)
			}
		}
	}
}
func (h *ContentHandler) Patch(w http.ResponseWriter, r *http.Request) {
	r, content, jsonContent, er1 := core.BuildMapAndCheckId[Content](w, r, h.Keys, h.Indexes)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &content)
		if !core.HasError(w, r, errors, er2, h.Error, jsonContent, h.Log, h.Resource, h.Action.Patch) {
			res, err := h.service.Patch(r.Context(), jsonContent)
			if err != nil {
				h.Error(r.Context(), err.Error())
				h.Log(r.Context(), h.Resource, h.Action.Patch, false, err.Error())
				http.Error(w, core.InternalServerError, http.StatusInternalServerError)
				return
			}

			if res > 0 {
				h.Log(r.Context(), h.Resource, h.Action.Patch, true, fmt.Sprintf("%s '%s' '%s'", h.Action.Patch, content.Id, content.Lang))
				core.JSON(w, http.StatusOK, jsonContent)
			} else if res == 0 {
				h.Log(r.Context(), h.Resource, h.Action.Patch, false, fmt.Sprintf("not found '%s' '%s'", content.Id, content.Lang))
				core.JSON(w, http.StatusNotFound, res)
			} else {
				h.Log(r.Context(), h.Resource, h.Action.Patch, false, fmt.Sprintf("conflict '%s' '%s'", content.Id, content.Lang))
				core.JSON(w, http.StatusConflict, res)
			}
		}
	}
}
func (h *ContentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, er1 := core.GetRequiredString(w, r, 1)
	lang, er2 := core.GetRequiredString(w, r)
	if er1 == nil && er2 == nil {
		res, err := h.service.Delete(r.Context(), id, lang)
		if err != nil {
			h.Error(r.Context(), err.Error())
			h.Log(r.Context(), h.Resource, h.Action.Delete, false, err.Error())
			http.Error(w, core.InternalServerError, http.StatusInternalServerError)
			return
		}

		if res > 0 {
			h.Log(r.Context(), h.Resource, h.Action.Delete, true, fmt.Sprintf("%s '%s' '%s'", h.Action.Delete, id, lang))
			core.JSON(w, http.StatusOK, res)
		} else if res == 0 {
			h.Log(r.Context(), h.Resource, h.Action.Delete, false, fmt.Sprintf("not found '%s' '%s'", id, lang))
			core.JSON(w, http.StatusNotFound, res)
		} else {
			h.Log(r.Context(), h.Resource, h.Action.Delete, false, fmt.Sprintf("conflict '%s' '%s'", id, lang))
			core.JSON(w, http.StatusConflict, res)
		}
	}
}
func (h *ContentHandler) Search(w http.ResponseWriter, r *http.Request) {
	filter := ContentFilter{Filter: &search.Filter{}}
	err := search.Decode(r, &filter, h.ParamIndex, h.FilterIndex)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	offset := search.GetOffset(filter.Limit, filter.Page)
	contents, total, err := h.service.Search(r.Context(), &filter, filter.Limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	core.JSON(w, http.StatusOK, &search.Result{List: &contents, Total: total})
}
