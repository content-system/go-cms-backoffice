package app

import (
	"context"
	"database/sql"

	"github.com/core-go/authentication"
	ah "github.com/core-go/authentication/handler"
	"github.com/core-go/authentication/mock"
	as "github.com/core-go/authentication/sql"
	"github.com/core-go/core/authorization"
	"github.com/core-go/core/code"
	"github.com/core-go/core/health"
	hs "github.com/core-go/core/health/sql"
	se "github.com/core-go/core/settings"
	"github.com/core-go/core/shortid"
	ur "github.com/core-go/core/user"
	log "github.com/core-go/log/zap"
	sec "github.com/core-go/security"
	"github.com/core-go/security/jwt"
	ss "github.com/core-go/security/sql"
	q "github.com/core-go/sql"
	sa "github.com/core-go/sql/action"
	"github.com/core-go/sql/template"
	"github.com/core-go/sql/template/xml"

	a "go-service/internal/article"
	"go-service/internal/audit-log"
	c "go-service/internal/contact"
	co "go-service/internal/content"
	j "go-service/internal/job"
	r "go-service/internal/role"
	u "go-service/internal/user"
	p "go-service/pkg/privilege"
)

type ApplicationContext struct {
	SkipSecurity         bool
	Health               *health.Handler
	Authorization        *authorization.Handler
	AuthorizationChecker *sec.AuthorizationChecker
	Authorizer           *sec.Authorizer
	Authentication       *ah.AuthenticationHandler
	Privileges           *ah.PrivilegesHandler
	Privilege            *p.PrivilegesHandler
	Code                 *code.Handler
	Roles                *code.Handler
	Role                 r.RoleTransport
	User                 u.UserTransport
	AuditLog             *audit.AuditLogHandler
	Settings             *se.Handler
	Content              co.ContentTransport
	Article              a.ArticleTransport
	Job                  j.JobTransport
	Contact              c.ContactTransport
}

func NewApp(ctx context.Context, cfg Config) (*ApplicationContext, error) {
	db, er0 := sql.Open(cfg.DB.Driver, cfg.DB.DataSourceName)
	if er0 != nil {
		return nil, er0
	}
	sqlHealthChecker := hs.NewHealthChecker(db)
	var healthHandler *health.Handler

	logError := log.LogError
	generateId := shortid.Generate
	var writeLog func(ctx context.Context, resource string, action string, success bool, desc string) error

	if cfg.AuditLog.Log {
		auditLogDB, er1 := sql.Open(cfg.AuditLog.DB.Driver, cfg.AuditLog.DB.DataSourceName)
		if er1 != nil {
			return nil, er1
		}
		logWriter := sa.NewActionLogWriter(auditLogDB, "audit_logs", cfg.AuditLog.Config, cfg.AuditLog.Schema, generateId)
		writeLog = logWriter.Write
		auditLogHealthChecker := hs.NewSqlHealthChecker(auditLogDB, "audit_logs")
		healthHandler = health.NewHandler(sqlHealthChecker, auditLogHealthChecker)
	} else {
		healthHandler = health.NewHandler(sqlHealthChecker)
	}
	buildParam := q.GetBuild(db)
	sqlPrivilegeLoader := ss.NewPrivilegeLoader(db, cfg.Sql.PermissionsByUser)

	userId := cfg.Tracking.User
	tokenPort := jwt.NewTokenAdapter()
	authorizationHandler := authorization.NewHandler(tokenPort.GetAndVerifyToken, cfg.Auth.Token.Secret)
	authorizationChecker := sec.NewAuthorizationChecker(tokenPort.GetAndVerifyToken, cfg.Auth.Token.Secret, userId)
	authorizer := sec.NewAuthorizer(sqlPrivilegeLoader.Privilege, true, userId)

	authStatus := auth.InitStatus(cfg.Auth.Status)
	ldapAuthenticator, er2 := mock.NewDAPAuthenticatorByConfig(cfg.Ldap, authStatus)
	if er2 != nil {
		return nil, er2
	}
	userPort, er3 := as.NewUserAdapter(db, cfg.Auth.Query, cfg.Auth.DB, cfg.Auth.UserStatus)
	if er3 != nil {
		return nil, er3
	}
	privilegePort, er4 := as.NewSqlPrivilegesAdapter(db, cfg.Sql.PrivilegesByUser, 1, true)
	if er4 != nil {
		return nil, er4
	}
	authenticator := auth.NewBasicAuthenticator(authStatus, ldapAuthenticator.Authenticate, userPort, tokenPort.GenerateToken, cfg.Auth.Token, cfg.Auth.Payload, privilegePort.Load)
	authenticationHandler := ah.NewAuthenticationHandler(authenticator.Authenticate, authStatus.Error, authStatus.Timeout, logError, writeLog)

	privilegeReader, er5 := as.NewPrivilegesReader(db, cfg.Sql.Privileges)
	if er5 != nil {
		return nil, er5
	}
	privilegesHandler := ah.NewPrivilegesHandler(privilegeReader.Privileges)
	privilegeHandler := p.NewPrivilegesHandler(privilegeReader.Privileges, privilegePort.Load, logError, "userId")

	// codeLoader := code.NewDynamicSqlCodeLoader(db, "select code, name, status as text from codeMaster where master = ? and status = 'A'", 1)
	codeLoader, err := code.NewSqlCodeLoader(db, "code_masters", cfg.Code.Loader)
	if err != nil {
		return nil, err
	}
	codeHandler := code.NewCodeHandlerByConfig(codeLoader.Load, cfg.Code.Handler, logError)

	templates, err := template.LoadTemplates(xml.Trim, "configs/query.xml")
	if err != nil {
		return nil, err
	}
	// rolesLoader, err := code.NewDynamicSqlCodeLoader(db, "select roleName as name, roleId as id from roles where status = 'A'", 0)

	rolesLoader, err := code.NewSqlCodeLoader(db, "roles", cfg.Role.Loader)
	if err != nil {
		return nil, err
	}
	rolesHandler := code.NewCodeHandlerByConfig(rolesLoader.Load, cfg.Role.Handler, logError)

	roleHandler, err := r.NewRoleTransport(db, logError, templates, cfg.Tracking, writeLog, cfg.Action)
	if err != nil {
		return nil, err
	}

	userHandler, err := u.NewUserTransport(db, logError, templates, cfg.Tracking, writeLog, cfg.Action)
	if err != nil {
		return nil, err
	}

	contentHandler, err := co.NewContentTransport(db, logError, writeLog, cfg.Action)
	if err != nil {
		return nil, err
	}

	articleHandler, err := a.NewArticleTransport(db, logError, writeLog, cfg.Action)
	if err != nil {
		return nil, err
	}
	jobHandler, err := j.NewJobTransport(db, logError, writeLog, cfg.Action)
	if err != nil {
		return nil, err
	}

	contactHandler, err := c.NewContactTransport(db, logError, writeLog, cfg.Action)
	if err != nil {
		return nil, err
	}

	reportDB, er8 := sql.Open(cfg.AuditLog.DB.Driver, cfg.AuditLog.DB.DataSourceName)
	if er8 != nil {
		return nil, er8
	}
	userQuery := ur.NewUserAdapter(db, "select user_id, display_name, email, phone, image_url from users where user_id ")

	auditLogQuery, er9 := audit.NewAuditLogQuery(reportDB, templates, userQuery.Query)
	if er9 != nil {
		return nil, er9
	}
	auditLogHandler := audit.NewAuditLogHandler(auditLogQuery, logError)

	settingsHandler := se.NewSettingsHandler(logError, writeLog, db, "users", buildParam, "userId", "user_id", "dateformat", "language")

	app := &ApplicationContext{
		Health:               healthHandler,
		SkipSecurity:         cfg.SecuritySkip,
		Authorization:        authorizationHandler,
		AuthorizationChecker: authorizationChecker,
		Authorizer:           authorizer,
		Authentication:       authenticationHandler,
		Privileges:           privilegesHandler,
		Privilege:            privilegeHandler,
		Code:                 codeHandler,
		Roles:                rolesHandler,
		Role:                 roleHandler,
		User:                 userHandler,
		AuditLog:             auditLogHandler,
		Settings:             settingsHandler,
		Content:              contentHandler,
		Article:              articleHandler,
		Job:                  jobHandler,
		Contact:              contactHandler,
	}
	return app, nil
}
