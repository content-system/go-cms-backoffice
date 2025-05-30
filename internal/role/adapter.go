package role

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	q "github.com/core-go/sql"
)

const ActionNone int32 = 0

type userRole struct {
	UserId string `json:"userId,omitempty" gorm:"column:user_id;primary_key" bson:"_id,omitempty" validate:"required,max=20,code"`
	RoleId string `json:"roleId,omitempty" gorm:"column:role_id;primary_key" bson:"_id,omitempty" dynamodbav:"roleId,omitempty" firestore:"roleId,omitempty" validate:"max=40"`
}

type roleModule struct {
	RoleId      string `json:"roleId,omitempty" gorm:"column:role_id" bson:"roleId,omitempty" dynamodbav:"roleId,omitempty" firestore:"roleId,omitempty" validate:"required"`
	ModuleId    string `json:"moduleId,omitempty" gorm:"column:module_id" bson:"moduleId,omitempty" dynamodbav:"moduleId,omitempty" firestore:"moduleId,omitempty" validate:"required"`
	Permissions int32  `json:"permissions,omitempty" gorm:"column:permissions" bson:"permissions,omitempty" dynamodbav:"permissions,omitempty" firestore:"permissions,omitempty" validate:"required"`
}

type RoleAdapter struct {
	db            *sql.DB
	Driver        string
	BuildParam    func(int) string
	keys          []string
	jsonColumnMap map[string]string
	Map           map[string]int
	Schema        *q.Schema
	ModuleMap     map[string]int
	ModuleSchema  *q.Schema
	UserSchema    *q.Schema
}

func NewRoleAdapter(db *sql.DB) (*RoleAdapter, error) {
	roleMap, roleSchema, jsonColumnMap, keys, _, _, buildParam, driver, err := q.Init(reflect.TypeOf(Role{}), db)
	if err != nil {
		return nil, err
	}
	userRoleSchema := q.CreateSchema(reflect.TypeOf(userRole{}))

	moduleType := reflect.TypeOf(roleModule{})
	roleModuleSchema := q.CreateSchema(moduleType)
	moduleMap, err := q.GetColumnIndexes(moduleType)

	return &RoleAdapter{
			db:            db,
			Driver:        driver,
			BuildParam:    buildParam,
			Map:           roleMap,
			Schema:        roleSchema,
			jsonColumnMap: jsonColumnMap,
			keys:          keys,
			ModuleMap:     moduleMap,
			ModuleSchema:  roleModuleSchema,
			UserSchema:    userRoleSchema,
		},
		err
}

func (s *RoleAdapter) Load(ctx context.Context, roleId string) (*Role, error) {
	var roles []Role
	query1 := fmt.Sprintf("select * from roles where role_id = %s", s.BuildParam(1))
	er1 := q.Query(ctx, s.db, s.Map, &roles, query1, roleId)
	if er1 != nil {
		return nil, er1
	}
	if len(roles) == 0 {
		return nil, nil
	}
	role := roles[0]
	var modules []roleModule
	query2 := fmt.Sprintf(`select module_id, permissions from role_modules where role_id = %s`, s.BuildParam(1))
	er2 := q.Query(ctx, s.db, s.ModuleMap, &modules, query2, roleId)
	if er2 != nil {
		return nil, er2
	}
	privileges := make([]string, 0)
	if len(modules) > 0 {
		for _, module := range modules {
			id := module.ModuleId
			if module.Permissions != 0 {
				id = module.ModuleId + " " + fmt.Sprintf("%X", module.Permissions)
			}
			privileges = append(privileges, id)
		}
	}

	role.Privileges = privileges
	return &role, nil
}

func buildModules(roleId string, privileges []string) ([]roleModule, error) {
	if privileges == nil || len(privileges) <= 0 {
		return nil, nil
	}
	modules := make([]roleModule, 0)
	for _, p := range privileges {
		m := toModules(p)
		m.RoleId = roleId
		modules = append(modules, m)
	}
	return modules, nil
}
func toModules(menu string) roleModule {
	s := strings.Split(menu, " ")
	permission := ActionNone
	if len(s) >= 2 {
		i, err := strconv.ParseInt(s[1], 16, 64)
		if err == nil {
			permission = int32(i)
		}
	}
	p := roleModule{ModuleId: s[0], Permissions: permission}
	return p
}
func (s *RoleAdapter) Create(ctx context.Context, role *Role) (int64, error) {
	modules, er1 := buildModules(role.RoleId, role.Privileges)
	if er1 != nil {
		return 0, er1
	}
	sts := q.NewStatements(true)
	sts.Add(q.BuildToInsert("roles", role, s.BuildParam, s.Schema))
	if modules != nil {
		query, args, er2 := q.BuildToInsertBatch("role_modules", modules, s.Driver, s.ModuleSchema)
		if er2 != nil {
			return 0, er2
		}
		sts.Add(query, args)
	}

	return sts.Exec(ctx, s.db)
}
func (s *RoleAdapter) Update(ctx context.Context, role *Role) (int64, error) {
	modules, err := buildModules(role.RoleId, role.Privileges)
	if err != nil {
		return 0, err
	}
	sts := q.NewStatements(true)
	sts.Add(q.BuildToUpdate("roles", role, s.BuildParam, s.Schema))

	deleteModules := fmt.Sprintf("delete from role_modules where role_id = %s", s.BuildParam(1))
	sts.Add(deleteModules, []interface{}{role.RoleId})
	if modules != nil {
		query, args, er2 := q.BuildToInsertBatch("role_modules", modules, s.Driver, s.ModuleSchema)
		if er2 != nil {
			return 0, er2
		}
		sts.Add(query, args)
	}

	return sts.Exec(ctx, s.db)
}

func (s *RoleAdapter) Patch(ctx context.Context, role map[string]interface{}) (int64, error) {
	objId, ok := role["roleId"]
	if !ok {
		return -1, errors.New("roleId must be in payload")
	}
	roleId, ok2 := objId.(string)
	if !ok2 {
		return -1, errors.New("roleId must be a string")
	}
	var privileges []string
	var ok4 bool
	objPrivileges, ok3 := role["privileges"]
	if ok3 {
		privileges, ok4 = objPrivileges.([]string)
	}
	var sts q.Statements
	if ok4 && len(role) > 2 || !ok4 && len(role) > 1 {
		sts = q.NewStatements(true)
		columnMap := q.JSONToColumns(role, s.jsonColumnMap)
		sts.Add(q.BuildToPatch("roles", columnMap, s.keys, s.BuildParam))
	} else {
		sts = q.NewStatements(false)
	}

	deleteModules := fmt.Sprintf("delete from role_modules where role_id = %s", s.BuildParam(1))
	sts.Add(deleteModules, []interface{}{roleId})

	if ok4 {
		modules, err := buildModules(roleId, privileges)
		if err != nil {
			return -1, err
		}
		query, args, er2 := q.BuildToInsertBatch("role_modules", modules, s.Driver, s.ModuleSchema)
		if er2 != nil {
			return -1, er2
		}
		sts.Add(query, args)
	}
	return sts.Exec(ctx, s.db)
}

func (s *RoleAdapter) Delete(ctx context.Context, id string) (int64, error) {
	exist, er0 := q.Exist(ctx, s.db, fmt.Sprintf("select user_id from user_roles where role_id = %s limit 1", s.BuildParam(1)), id)
	if exist || er0 != nil {
		return -1, er0
	}

	sts := q.NewStatements(false)

	deleteModules := fmt.Sprintf("delete from role_modules where role_id = %s", s.BuildParam(1))
	sts.Add(deleteModules, []interface{}{id})

	deleteRole := fmt.Sprintf("delete from roles where role_id = %s", s.BuildParam(1))
	sts.Add(deleteRole, []interface{}{id})

	return sts.Exec(ctx, s.db)
}

func (s *RoleAdapter) AssignRole(ctx context.Context, roleId string, users []string) (int64, error) {
	modules := make([]userRole, 0)
	for _, u := range users {
		modules = append(modules, userRole{UserId: u, RoleId: roleId})
	}
	sts := q.NewStatements(false)

	deleteModules := fmt.Sprintf("delete from user_roles where role_id = %s", s.BuildParam(1))
	sts.Add(deleteModules, []interface{}{roleId})
	if modules != nil {
		query, args, er2 := q.BuildToInsertBatch("user_roles", modules, s.Driver, s.UserSchema)
		if er2 != nil {
			return 0, er2
		}
		sts.Add(query, args)
	}

	return sts.Exec(ctx, s.db)
}
