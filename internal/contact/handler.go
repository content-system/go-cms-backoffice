package contact

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/core-go/core"
	"github.com/core-go/search"
)

func NewContactHandler(service ContactService, logError core.Log, validate core.Validate[*Contact], writeLog core.WriteLog, action *core.ActionConfig) *ContactHandler {
	contactType := reflect.TypeOf(Contact{})
	parameters := search.CreateParameters(reflect.TypeOf(ContactFilter{}), contactType)
	attributes := core.CreateAttributes(contactType, logError, action, writeLog)
	return &ContactHandler{service: service, Validate: validate, Attributes: attributes, Parameters: parameters}
}

type ContactHandler struct {
	service  ContactService
	Validate core.Validate[*Contact]
	*core.Attributes
	*search.Parameters
}

func (h *ContactHandler) Load(w http.ResponseWriter, r *http.Request) {
	id, err := core.GetRequiredString(w, r)
	if err == nil {
		contact, err := h.service.Load(r.Context(), id)
		if err != nil {
			h.Error(r.Context(), fmt.Sprintf("Error to get contact '%s': %s", id, err.Error()))
			http.Error(w, core.InternalServerError, http.StatusInternalServerError)
			return
		}
		core.JSON(w, core.IsFound(contact), contact)
	}
}
func (h *ContactHandler) Create(w http.ResponseWriter, r *http.Request) {
	contact, er1 := core.Decode[Contact](w, r)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &contact)
		if !core.HasError(w, r, errors, er2, h.Error, &contact, h.Log, h.Resource, h.Action.Create) {
			res, er3 := h.service.Create(r.Context(), &contact)
			if er3 != nil {
				h.Error(r.Context(), er3.Error())
				h.Log(r.Context(), h.Resource, h.Action.Update, false, er3.Error())
				http.Error(w, core.InternalServerError, http.StatusInternalServerError)
				return
			}

			if res > 0 {
				h.Log(r.Context(), h.Resource, h.Action.Update, true, fmt.Sprintf("create '%s'", contact.Id))
				core.JSON(w, http.StatusCreated, contact)
			} else {
				h.Log(r.Context(), h.Resource, h.Action.Update, false, fmt.Sprintf("conflict '%s'", contact.Id))
				core.JSON(w, http.StatusConflict, res)
			}
		}
	}
}
func (h *ContactHandler) Update(w http.ResponseWriter, r *http.Request) {
	contact, er1 := core.DecodeAndCheckId[Contact](w, r, h.Keys, h.Indexes)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &contact)
		if !core.HasError(w, r, errors, er2, h.Error, &contact, h.Log, h.Resource, h.Action.Update) {
			res, err := h.service.Update(r.Context(), &contact)
			if err != nil {
				h.Error(r.Context(), err.Error())
				h.Log(r.Context(), h.Resource, h.Action.Update, false, err.Error())
				http.Error(w, core.InternalServerError, http.StatusInternalServerError)
				return
			}

			if res > 0 {
				h.Log(r.Context(), h.Resource, h.Action.Update, true, fmt.Sprintf("%s '%s'", h.Action.Update, contact.Id))
				core.JSON(w, http.StatusOK, contact)
			} else if res == 0 {
				h.Log(r.Context(), h.Resource, h.Action.Update, false, fmt.Sprintf("not found '%s'", contact.Id))
				core.JSON(w, http.StatusNotFound, res)
			} else {
				h.Log(r.Context(), h.Resource, h.Action.Update, false, fmt.Sprintf("conflict '%s'", contact.Id))
				core.JSON(w, http.StatusConflict, res)
			}
		}
	}
}
func (h *ContactHandler) Patch(w http.ResponseWriter, r *http.Request) {
	r, contact, jsonContact, er1 := core.BuildMapAndCheckId[Contact](w, r, h.Keys, h.Indexes)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &contact)
		if !core.HasError(w, r, errors, er2, h.Error, jsonContact, h.Log, h.Resource, h.Action.Patch) {
			res, err := h.service.Patch(r.Context(), jsonContact)
			if err != nil {
				h.Error(r.Context(), err.Error())
				h.Log(r.Context(), h.Resource, h.Action.Patch, false, err.Error())
				http.Error(w, core.InternalServerError, http.StatusInternalServerError)
				return
			}

			if res > 0 {
				h.Log(r.Context(), h.Resource, h.Action.Patch, true, fmt.Sprintf("%s '%s'", h.Action.Patch, contact.Id))
				core.JSON(w, http.StatusOK, jsonContact)
			} else if res == 0 {
				h.Log(r.Context(), h.Resource, h.Action.Patch, false, fmt.Sprintf("not found '%s'", contact.Id))
				core.JSON(w, http.StatusNotFound, res)
			} else {
				h.Log(r.Context(), h.Resource, h.Action.Patch, false, fmt.Sprintf("conflict '%s'", contact.Id))
				core.JSON(w, http.StatusConflict, res)
			}
		}
	}
}
func (h *ContactHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
func (h *ContactHandler) Search(w http.ResponseWriter, r *http.Request) {
	filter := ContactFilter{Filter: &search.Filter{}}
	err := search.Decode(r, &filter, h.ParamIndex, h.FilterIndex)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	offset := search.GetOffset(filter.Limit, filter.Page)
	contacts, total, err := h.service.Search(r.Context(), &filter, filter.Limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	core.JSON(w, http.StatusOK, &search.Result{List: &contacts, Total: total})
}
