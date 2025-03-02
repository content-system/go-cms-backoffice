package job

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/core-go/core"
	"github.com/core-go/search"
)

func NewJobHandler(service JobService, logError core.Log, validate core.Validate[*Job], action *core.ActionConfig) *JobHandler {
	jobType := reflect.TypeOf(Job{})
	parameters := search.CreateParameters(reflect.TypeOf(JobFilter{}), jobType)
	attributes := core.CreateAttributes(jobType, logError, action)
	return &JobHandler{service: service, Validate: validate, Attributes: attributes, Parameters: parameters}
}

type JobHandler struct {
	service  JobService
	Validate core.Validate[*Job]
	*core.Attributes
	*search.Parameters
}

func (h *JobHandler) Load(w http.ResponseWriter, r *http.Request) {
	id, err := core.GetRequiredString(w, r)
	if err == nil {
		job, err := h.service.Load(r.Context(), id)
		if err != nil {
			h.Error(r.Context(), fmt.Sprintf("Error to get job '%s': %s", id, err.Error()))
			http.Error(w, core.InternalServerError, http.StatusInternalServerError)
			return
		}
		core.JSON(w, core.IsFound(job), job)
	}
}
func (h *JobHandler) Create(w http.ResponseWriter, r *http.Request) {
	job, er1 := core.Decode[Job](w, r)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &job)
		if !core.HasError(w, r, errors, er2, h.Error, &job, h.Log, h.Resource, h.Action.Create) {
			res, er3 := h.service.Create(r.Context(), &job)
			if er3 != nil {
				h.Error(r.Context(), er3.Error())
				h.Log(r.Context(), h.Resource, h.Action.Update, false, er3.Error())
				http.Error(w, core.InternalServerError, http.StatusInternalServerError)
				return
			}

			if res > 0 {
				h.Log(r.Context(), h.Resource, h.Action.Create, true, fmt.Sprintf("%s '%s'", h.Action.Create, job.Id))
				core.JSON(w, http.StatusCreated, job)
			} else {
				h.Log(r.Context(), h.Resource, h.Action.Create, false, fmt.Sprintf("conflict '%s'", job.Id))
				core.JSON(w, http.StatusConflict, res)
			}
		}
	}
}
func (h *JobHandler) Update(w http.ResponseWriter, r *http.Request) {
	job, er1 := core.DecodeAndCheckId[Job](w, r, h.Keys, h.Indexes)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &job)
		if !core.HasError(w, r, errors, er2, h.Error, &job, h.Log, h.Resource, h.Action.Update) {
			res, err := h.service.Update(r.Context(), &job)
			if err != nil {
				h.Error(r.Context(), err.Error())
				h.Log(r.Context(), h.Resource, h.Action.Update, false, err.Error())
				http.Error(w, core.InternalServerError, http.StatusInternalServerError)
				return
			}

			if res > 0 {
				h.Log(r.Context(), h.Resource, h.Action.Update, true, fmt.Sprintf("%s '%s'", h.Action.Update, job.Id))
				core.JSON(w, http.StatusOK, job)
			} else if res == 0 {
				h.Log(r.Context(), h.Resource, h.Action.Update, false, fmt.Sprintf("not found '%s'", job.Id))
				core.JSON(w, http.StatusNotFound, res)
			} else {
				h.Log(r.Context(), h.Resource, h.Action.Update, false, fmt.Sprintf("conflict '%s'", job.Id))
				core.JSON(w, http.StatusConflict, res)
			}
		}
	}
}
func (h *JobHandler) Patch(w http.ResponseWriter, r *http.Request) {
	r, job, jsonJob, er1 := core.BuildMapAndCheckId[Job](w, r, h.Keys, h.Indexes)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &job)
		if !core.HasError(w, r, errors, er2, h.Error, jsonJob, h.Log, h.Resource, h.Action.Patch) {
			res, err := h.service.Patch(r.Context(), jsonJob)
			if err != nil {
				h.Error(r.Context(), err.Error())
				h.Log(r.Context(), h.Resource, h.Action.Patch, false, err.Error())
				http.Error(w, core.InternalServerError, http.StatusInternalServerError)
				return
			}

			if res > 0 {
				h.Log(r.Context(), h.Resource, h.Action.Patch, true, fmt.Sprintf("%s '%s'", h.Action.Patch, job.Id))
				core.JSON(w, http.StatusOK, jsonJob)
			} else if res == 0 {
				h.Log(r.Context(), h.Resource, h.Action.Patch, false, fmt.Sprintf("not found '%s'", job.Id))
				core.JSON(w, http.StatusNotFound, res)
			} else {
				h.Log(r.Context(), h.Resource, h.Action.Patch, false, fmt.Sprintf("conflict '%s'", job.Id))
				core.JSON(w, http.StatusConflict, res)
			}
		}
	}
}
func (h *JobHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
func (h *JobHandler) Search(w http.ResponseWriter, r *http.Request) {
	filter := JobFilter{Filter: &search.Filter{}}
	err := search.Decode(r, &filter, h.ParamIndex, h.FilterIndex)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	offset := search.GetOffset(filter.Limit, filter.Page)
	jobs, total, err := h.service.Search(r.Context(), &filter, filter.Limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	core.JSON(w, http.StatusOK, &search.Result{List: &jobs, Total: total})
}
