package app

import (
	"context"
	"net/http"

	c "github.com/core-go/core/constants"
	m "github.com/core-go/core/mux"
	s "github.com/core-go/core/security"
	"github.com/gorilla/mux"
)

const (
	role      = "role"
	user      = "user"
	audit_log = "audit_log"
	category  = "category"
	content   = "content"
	article   = "article"
	job       = "job"
	contact   = "contact"
)

func Route(r *mux.Router, ctx context.Context, conf Config) error {
	app, err := NewApp(ctx, conf)
	if err != nil {
		return err
	}
	r.Use(app.Authorization.HandleAuthorization)
	sec := &s.SecurityConfig{SecuritySkip: conf.SecuritySkip, Check: app.AuthorizationChecker.Check, Authorize: app.Authorizer.Authorize}

	Handle(r, "/health", app.Health.Check, c.GET)
	Handle(r, "/authenticate", app.Authentication.Authenticate, c.POST)

	r.Handle("/code/{code}", app.AuthorizationChecker.Check(http.HandlerFunc(app.Code.Load))).Methods(c.GET)
	r.Handle("/settings", app.AuthorizationChecker.Check(http.HandlerFunc(app.Settings.Save))).Methods(c.PATCH)

	Handle(r, "/my-privileges", app.Privilege.GetPrivileges, c.GET)

	HandleWithSecurity(sec, r, "/privileges", app.Privileges.All, role, c.ActionRead, c.GET)
	roles := r.PathPrefix("/roles").Subrouter()
	HandleWithSecurity(sec, roles, "/search", app.Role.Search, role, c.ActionRead, c.POST, c.GET)
	HandleWithSecurity(sec, roles, "/{roleId}", app.Role.Load, role, c.ActionRead, c.GET)
	HandleWithSecurity(sec, roles, "", app.Role.Create, role, c.ActionWrite, c.POST)
	HandleWithSecurity(sec, roles, "/{roleId}", app.Role.Update, role, c.ActionWrite, c.PUT)
	HandleWithSecurity(sec, roles, "/{userId}", app.Role.Patch, user, c.ActionWrite, c.PATCH)
	HandleWithSecurity(sec, roles, "/{roleId}", app.Role.Delete, role, c.ActionWrite, c.DELETE)
	HandleWithSecurity(sec, roles, "/{roleId}/assign", app.Role.AssignRole, role, c.ActionWrite, c.PUT)

	HandleWithSecurity(sec, r, "/roles", app.Roles.Load, user, c.ActionRead, c.GET)
	users := r.PathPrefix("/users").Subrouter()
	HandleWithSecurity(sec, users, "", app.User.GetUserByRole, role, c.ActionRead, c.GET)

	m.HandleWithSecurity(users, "/search", app.User.Search, sec.Check, sec.Authorize, user, c.ActionRead, c.GET, c.POST)
	HandleWithSecurity(sec, users, "/{userId}", app.User.Load, user, c.ActionRead, c.GET)
	HandleWithSecurity(sec, users, "", app.User.Create, user, c.ActionWrite, c.POST)
	HandleWithSecurity(sec, users, "/{userId}", app.User.Update, user, c.ActionWrite, c.PUT)
	HandleWithSecurity(sec, users, "/{userId}", app.User.Patch, user, c.ActionWrite, c.PATCH)
	HandleWithSecurity(sec, users, "/{userId}", app.User.Delete, user, c.ActionWrite, c.DELETE)

	categories := r.PathPrefix("/categories").Subrouter()
	HandleWithSecurity(sec, categories, "/search", app.Category.Search, category, c.ActionRead, c.GET, c.POST)
	HandleWithSecurity(sec, categories, "/{id}", app.Category.Load, category, c.ActionRead, c.GET)
	HandleWithSecurity(sec, categories, "", app.Category.Create, category, c.ActionWrite, c.POST)
	HandleWithSecurity(sec, categories, "/{id}", app.Category.Update, category, c.ActionWrite, c.PUT)
	HandleWithSecurity(sec, categories, "/{id}", app.Category.Patch, category, c.ActionWrite, c.PATCH)
	HandleWithSecurity(sec, categories, "/{id}", app.Category.Delete, category, c.ActionWrite, c.DELETE)

	contents := r.PathPrefix("/contents").Subrouter()
	HandleWithSecurity(sec, contents, "/search", app.Content.Search, content, c.ActionRead, c.GET, c.POST)
	HandleWithSecurity(sec, contents, "/{id}/{lang}", app.Content.Load, content, c.ActionRead, c.GET)
	HandleWithSecurity(sec, contents, "", app.Content.Create, content, c.ActionWrite, c.POST)
	HandleWithSecurity(sec, contents, "/{id}/{lang}", app.Content.Update, content, c.ActionWrite, c.PUT)
	HandleWithSecurity(sec, contents, "/{id}/{lang}", app.Content.Patch, content, c.ActionWrite, c.PATCH)
	HandleWithSecurity(sec, contents, "/{id}/{lang}", app.Content.Delete, content, c.ActionWrite, c.DELETE)

	articles := r.PathPrefix("/articles").Subrouter()
	HandleWithSecurity(sec, articles, "/search", app.Article.Search, article, c.ActionRead, c.GET, c.POST)
	HandleWithSecurity(sec, articles, "/{id}", app.Article.Load, article, c.ActionRead, c.GET)
	HandleWithSecurity(sec, articles, "", app.Article.Create, article, c.ActionWrite, c.POST)
	HandleWithSecurity(sec, articles, "/{id}", app.Article.Update, article, c.ActionWrite, c.PUT)
	HandleWithSecurity(sec, articles, "/{id}", app.Article.Patch, article, c.ActionWrite, c.PATCH)
	HandleWithSecurity(sec, articles, "/{id}", app.Article.Delete, article, c.ActionWrite, c.DELETE)

	jobs := r.PathPrefix("/jobs").Subrouter()
	HandleWithSecurity(sec, jobs, "/search", app.Job.Search, job, c.ActionRead, c.GET, c.POST)
	HandleWithSecurity(sec, jobs, "/{id}", app.Job.Load, job, c.ActionRead, c.GET)
	HandleWithSecurity(sec, jobs, "", app.Job.Create, job, c.ActionWrite, c.POST)
	HandleWithSecurity(sec, jobs, "/{id}", app.Job.Update, job, c.ActionWrite, c.PUT)
	HandleWithSecurity(sec, jobs, "/{id}", app.Job.Patch, job, c.ActionWrite, c.PATCH)
	HandleWithSecurity(sec, jobs, "/{id}", app.Job.Delete, job, c.ActionWrite, c.DELETE)

	contacts := r.PathPrefix("/contacts").Subrouter()
	HandleWithSecurity(sec, contacts, "/search", app.Contact.Search, contact, c.ActionRead, c.GET, c.POST)
	HandleWithSecurity(sec, contacts, "/{contactId}", app.Contact.Load, contact, c.ActionRead, c.GET)
	HandleWithSecurity(sec, contacts, "", app.Contact.Create, contact, c.ActionWrite, c.POST)
	HandleWithSecurity(sec, contacts, "/{contactId}", app.Contact.Update, contact, c.ActionWrite, c.PUT)
	HandleWithSecurity(sec, contacts, "/{contactId}", app.Contact.Patch, contact, c.ActionWrite, c.PATCH)
	HandleWithSecurity(sec, contacts, "/{contactId}", app.Contact.Delete, contact, c.ActionWrite, c.DELETE)

	HandleWithSecurity(sec, r, "/audit-logs", app.AuditLog.Search, audit_log, c.ActionRead, c.GET, c.POST)
	HandleWithSecurity(sec, r, "/audit-logs/search", app.AuditLog.Search, audit_log, c.ActionRead, c.GET, c.POST)
	return nil
}

func Handle(r *mux.Router, path string, f func(http.ResponseWriter, *http.Request), methods ...string) *mux.Route {
	return r.HandleFunc(path, f).Methods(methods...)
}
func HandleWithSecurity(authorizer *s.SecurityConfig, r *mux.Router, path string, f func(http.ResponseWriter, *http.Request), menuId string, action int32, methods ...string) *mux.Route {
	finalHandler := http.HandlerFunc(f)
	if authorizer.SecuritySkip {
		return r.HandleFunc(path, finalHandler).Methods(methods...)
	}
	authorize := func(next http.Handler) http.Handler {
		return authorizer.Authorize(next, menuId, action)
	}
	return r.Handle(path, authorizer.Check(authorize(finalHandler))).Methods(methods...)
}
