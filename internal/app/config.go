package app

import (
	q "github.com/core-go/authentication/sql"
	"github.com/core-go/core"
	"github.com/core-go/core/audit"
	"github.com/core-go/core/builder"
	"github.com/core-go/core/code"
	"github.com/core-go/core/cors"
	"github.com/core-go/core/server"
	mid "github.com/core-go/log/middleware"
	"github.com/core-go/log/zap"
	sa "github.com/core-go/sql/action"
)

type Config struct {
	Server       server.ServerConfig    `mapstructure:"server"`
	Allow        cors.AllowConfig       `mapstructure:"allow"`
	SecuritySkip bool                   `mapstructure:"security_skip"`
	Template     bool                   `mapstructure:"template"`
	Auth         q.SqlAuthConfig        `mapstructure:"auth"`
	DB           DBConfig               `mapstructure:"db"`
	Log          log.Config             `mapstructure:"log"`
	MiddleWare   mid.LogConfig          `mapstructure:"middleware"`
	AutoRoleId   *bool                  `mapstructure:"auto_role_id"`
	AutoUserId   *bool                  `mapstructure:"auto_user_id"`
	Role         code.Config            `mapstructure:"role"`
	Code         code.Config            `mapstructure:"code"`
	AuditLog     sa.ActionLogConf       `mapstructure:"audit_log"`
	AuditClient  audit.AuditLogClient   `mapstructure:"audit_client"`
	Action       *core.ActionConfig     `mapstructure:"action"`
	Tracking     builder.TrackingConfig `mapstructure:"tracking"`
	Sql          SqlStatement           `mapstructure:"sql"`
}
type DBConfig struct {
	DataSourceName string `yaml:"data_source_name" mapstructure:"data_source_name" json:"dataSourceName,omitempty" gorm:"column:datasourcename" bson:"dataSourceName,omitempty" dynamodbav:"dataSourceName,omitempty" firestore:"dataSourceName,omitempty"`
	Driver         string `yaml:"driver" mapstructure:"driver" json:"driver,omitempty" gorm:"column:driver" bson:"driver,omitempty" dynamodbav:"driver,omitempty" firestore:"driver,omitempty"`
}
type SqlStatement struct {
	Privileges        string `mapstructure:"privileges"`
	PrivilegesByUser  string `mapstructure:"privileges_by_user"`
	PermissionsByUser string `mapstructure:"permissions_by_user"`
}
