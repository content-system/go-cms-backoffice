package contact

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	s "github.com/core-go/sql"
)

func NewContactAdapter(db *sql.DB, buildQuery func(*ContactFilter) (string, []interface{})) (*ContactAdapter, error) {
	parameters, err := s.CreateParameters(reflect.TypeOf(Contact{}), db)
	if err != nil {
		return nil, err
	}
	return &ContactAdapter{DB: db, Parameters: parameters, BuildQuery: buildQuery}, nil
}

type ContactAdapter struct {
	DB         *sql.DB
	BuildQuery func(*ContactFilter) (string, []interface{})
	*s.Parameters
}

func (r *ContactAdapter) All(ctx context.Context) ([]Contact, error) {
	query := `select * from contacts`
	var contacts []Contact
	err := s.Query(ctx, r.DB, r.Map, &contacts, query)
	return contacts, err
}

func (r *ContactAdapter) Load(ctx context.Context, id string) (*Contact, error) {
	var contacts []Contact
	query := fmt.Sprintf("select %s from contacts where id = %s limit 1", r.Fields, r.BuildParam(1))
	err := s.Query(ctx, r.DB, r.Map, &contacts, query, id)
	if err != nil {
		return nil, err
	}
	if len(contacts) > 0 {
		return &contacts[0], nil
	}
	return nil, nil
}

func (r *ContactAdapter) Create(ctx context.Context, contact *Contact) (int64, error) {
	query, args := s.BuildToInsert("contacts", contact, r.BuildParam, r.Schema)
	tx := s.GetTx(ctx, r.DB)
	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func (r *ContactAdapter) Update(ctx context.Context, contact *Contact) (int64, error) {
	query, args := s.BuildToUpdate("contacts", contact, r.BuildParam, r.Schema)
	tx := s.GetTx(ctx, r.DB)
	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func (r *ContactAdapter) Patch(ctx context.Context, contact map[string]interface{}) (int64, error) {
	colMap := s.JSONToColumns(contact, r.JsonColumnMap)
	query, args := s.BuildToPatch("contacts", colMap, r.Keys, r.BuildParam)
	tx := s.GetTx(ctx, r.DB)
	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func (r *ContactAdapter) Delete(ctx context.Context, id string) (int64, error) {
	query := fmt.Sprintf("delete from contacts where id = %s", r.BuildParam(1))
	tx := s.GetTx(ctx, r.DB)
	res, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func (r *ContactAdapter) Search(ctx context.Context, filter *ContactFilter, limit int64, offset int64) ([]Contact, int64, error) {
	var contacts []Contact
	if limit <= 0 {
		return contacts, 0, nil
	}
	query, params := r.BuildQuery(filter)
	pagingQuery := s.BuildPagingQuery(query, limit, offset)
	countQuery := s.BuildCountQuery(query)

	row := r.DB.QueryRowContext(ctx, countQuery, params...)
	if row.Err() != nil {
		return contacts, 0, row.Err()
	}
	var total int64
	err := row.Scan(&total)
	if err != nil || total == 0 {
		return contacts, total, err
	}

	err = s.Query(ctx, r.DB, r.Map, &contacts, pagingQuery, params...)
	return contacts, total, err
}

func BuildQuery(filter *ContactFilter) (string, []interface{}) {
	query := "select * from contacts"
	where, params := BuildFilter(filter)
	if len(where) > 0 {
		query = query + " where " + where
	}
	return query, params
}
func BuildFilter(filter *ContactFilter) (string, []interface{}) {
	buildParam := s.BuildDollarParam
	var where []string
	var params []interface{}
	i := 1
	if len(filter.Id) > 0 {
		params = append(params, filter.Id)
		where = append(where, fmt.Sprintf(`id = %s`, buildParam(i)))
		i++
	}
	if filter.SubmittedAt != nil {
		if filter.SubmittedAt.Min != nil {
			params = append(params, filter.SubmittedAt.Min)
			where = append(where, fmt.Sprintf(`submitted_at >= %s`, buildParam(i)))
			i++
		}
		if filter.SubmittedAt.Max != nil {
			params = append(params, filter.SubmittedAt.Max)
			where = append(where, fmt.Sprintf(`submitted_at <= %s`, buildParam(i)))
			i++
		}
	}
	if len(filter.Email) > 0 {
		q := filter.Email + "%"
		params = append(params, q)
		where = append(where, fmt.Sprintf(`email ilike %s`, buildParam(i)))
		i++
	}
	if len(filter.Phone) > 0 {
		q := filter.Phone + "%"
		params = append(params, q)
		where = append(where, fmt.Sprintf(`phone ilike %s`, buildParam(i)))
		i++
	}
	if len(filter.Country) > 0 {
		q := filter.Country + "%"
		params = append(params, q)
		where = append(where, fmt.Sprintf(`country ilike %s`, buildParam(i)))
		i++
	}
	if len(filter.Name) > 0 {
		q := "%" + filter.Name + "%"
		params = append(params, q)
		where = append(where, fmt.Sprintf(`title ilike %s`, buildParam(i)))
		i++
	}
	if len(where) > 0 {
		return strings.Join(where, " and "), params
	}
	return "", params
}
