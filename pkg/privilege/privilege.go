package privilege

import (
	"context"
	"encoding/json"
	"net/http"

	au "github.com/core-go/authentication"
)

const internalServerError = "Internal Server Error"

type PrivilegesHandler struct {
	all   func(ctx context.Context) ([]au.Privilege, error)
	load  func(ctx context.Context, id string) ([]au.Privilege, error)
	Error func(context.Context, string, ...map[string]interface{})
	Key   string
}

func NewPrivilegesHandler(all func(context.Context) ([]au.Privilege, error), load func(ctx context.Context, id string) ([]au.Privilege, error), logError func(context.Context, string, ...map[string]interface{}), opts ...string) *PrivilegesHandler {
	user := "userId"
	if len(opts) > 0 && len(opts[0]) > 0 {
		user = opts[0]
	}
	h := PrivilegesHandler{all: all, load: load, Error: logError, Key: user}
	return &h
}
func (c *PrivilegesHandler) GetPrivileges(w http.ResponseWriter, r *http.Request) {
	userId := FromContext(r.Context(), c.Key)
	var privileges []au.Privilege
	var err error
	if len(userId) > 0 {
		privileges, err = c.load(r.Context(), userId)
	} else {
		privileges, err = c.all(r.Context())
	}
	if err != nil {
		if c.Error != nil {
			c.Error(r.Context(), err.Error())
		}
		c.Error(r.Context(), "error to get privileges: "+err.Error())
		JSON(w, http.StatusInternalServerError, internalServerError)
	} else {
		JSON(w, http.StatusOK, privileges)
	}
}
func JSON(w http.ResponseWriter, code int, result interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(result)
	return err
}
func FromContext(ctx context.Context, key string) string {
	u := ctx.Value(key)
	if u != nil {
		v, ok := u.(string)
		if ok {
			return v
		}
	}
	return ""
}
