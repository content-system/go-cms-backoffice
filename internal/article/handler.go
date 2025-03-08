package article

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/core-go/core"
	"github.com/core-go/search"
)

func NewArticleHandler(service ArticleService, logError core.Log, validate core.Validate[*Article], writeLog core.WriteLog, action *core.ActionConfig) *ArticleHandler {
	articleType := reflect.TypeOf(Article{})
	parameters := search.CreateParameters(reflect.TypeOf(ArticleFilter{}), articleType)
	attributes := core.CreateAttributes(articleType, logError, writeLog, action)
	return &ArticleHandler{service: service, Validate: validate, Attributes: attributes, Parameters: parameters}
}

type ArticleHandler struct {
	service  ArticleService
	Validate core.Validate[*Article]
	*core.Attributes
	*search.Parameters
}

func (h *ArticleHandler) Load(w http.ResponseWriter, r *http.Request) {
	id, err := core.GetRequiredString(w, r)
	if err == nil {
		article, err := h.service.Load(r.Context(), id)
		if err != nil {
			h.Error(r.Context(), fmt.Sprintf("Error to get article '%s': %s", id, err.Error()))
			http.Error(w, core.InternalServerError, http.StatusInternalServerError)
			return
		}
		core.JSON(w, core.IsFound(article), article)
	}
}
func (h *ArticleHandler) Create(w http.ResponseWriter, r *http.Request) {
	article, er1 := core.Decode[Article](w, r)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &article)
		if !core.HasError(w, r, errors, er2, h.Error, &article, h.Log, h.Resource, h.Action.Create) {
			res, er3 := h.service.Create(r.Context(), &article)
			if er3 != nil {
				h.Error(r.Context(), er3.Error())
				h.Log(r.Context(), h.Resource, h.Action.Update, false, er3.Error())
				http.Error(w, core.InternalServerError, http.StatusInternalServerError)
				return
			}

			if res > 0 {
				h.Log(r.Context(), h.Resource, h.Action.Create, true, fmt.Sprintf("%s '%s'", h.Action.Create, article.Id))
				core.JSON(w, http.StatusCreated, article)
			} else {
				h.Log(r.Context(), h.Resource, h.Action.Create, false, fmt.Sprintf("conflict '%s'", article.Id))
				core.JSON(w, http.StatusConflict, res)
			}
		}
	}
}
func (h *ArticleHandler) Update(w http.ResponseWriter, r *http.Request) {
	article, er1 := core.DecodeAndCheckId[Article](w, r, h.Keys, h.Indexes)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &article)
		if !core.HasError(w, r, errors, er2, h.Error, &article, h.Log, h.Resource, h.Action.Update) {
			res, err := h.service.Update(r.Context(), &article)
			if err != nil {
				h.Error(r.Context(), err.Error())
				h.Log(r.Context(), h.Resource, h.Action.Update, false, err.Error())
				http.Error(w, core.InternalServerError, http.StatusInternalServerError)
				return
			}

			if res > 0 {
				h.Log(r.Context(), h.Resource, h.Action.Update, true, fmt.Sprintf("%s '%s'", h.Action.Update, article.Id))
				core.JSON(w, http.StatusOK, article)
			} else if res == 0 {
				h.Log(r.Context(), h.Resource, h.Action.Update, false, fmt.Sprintf("not found '%s'", article.Id))
				core.JSON(w, http.StatusNotFound, res)
			} else {
				h.Log(r.Context(), h.Resource, h.Action.Update, false, fmt.Sprintf("conflict '%s'", article.Id))
				core.JSON(w, http.StatusConflict, res)
			}
		}
	}
}
func (h *ArticleHandler) Patch(w http.ResponseWriter, r *http.Request) {
	r, article, jsonArticle, er1 := core.BuildMapAndCheckId[Article](w, r, h.Keys, h.Indexes)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &article)
		if !core.HasError(w, r, errors, er2, h.Error, jsonArticle, h.Log, h.Resource, h.Action.Patch) {
			res, err := h.service.Patch(r.Context(), jsonArticle)
			if err != nil {
				h.Error(r.Context(), err.Error())
				h.Log(r.Context(), h.Resource, h.Action.Patch, false, err.Error())
				http.Error(w, core.InternalServerError, http.StatusInternalServerError)
				return
			}

			if res > 0 {
				h.Log(r.Context(), h.Resource, h.Action.Patch, true, fmt.Sprintf("%s '%s'", h.Action.Patch, article.Id))
				core.JSON(w, http.StatusOK, jsonArticle)
			} else if res == 0 {
				h.Log(r.Context(), h.Resource, h.Action.Patch, false, fmt.Sprintf("not found '%s'", article.Id))
				core.JSON(w, http.StatusNotFound, res)
			} else {
				h.Log(r.Context(), h.Resource, h.Action.Patch, false, fmt.Sprintf("conflict '%s'", article.Id))
				core.JSON(w, http.StatusConflict, res)
			}
		}
	}
}
func (h *ArticleHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
func (h *ArticleHandler) Search(w http.ResponseWriter, r *http.Request) {
	filter := ArticleFilter{Filter: &search.Filter{}}
	err := search.Decode(r, &filter, h.ParamIndex, h.FilterIndex)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	offset := search.GetOffset(filter.Limit, filter.Page)
	articles, total, err := h.service.Search(r.Context(), &filter, filter.Limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	core.JSON(w, http.StatusOK, &search.Result{List: &articles, Total: total})
}
