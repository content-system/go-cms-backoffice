package category

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	s "github.com/core-go/sql"
)

func NewCategoryAdapter(db *sql.DB, buildQuery func(*CategoryFilter) (string, []interface{})) (*CategoryAdapter, error) {
	parameters, err := s.CreateParameters(reflect.TypeOf(Category{}), db)
	if err != nil {
		return nil, err
	}
	versionIndex := parameters.Map["version"]
	return &CategoryAdapter{DB: db, Parameters: parameters, VersionIndex: versionIndex, BuildQuery: buildQuery}, nil
}

type CategoryAdapter struct {
	DB         *sql.DB
	BuildQuery func(*CategoryFilter) (string, []interface{})
	*s.Parameters
	VersionIndex int // index of field Version in the Category struct
}

func (r *CategoryAdapter) All(ctx context.Context) ([]Category, error) {
	query := `select * from categories`
	var categories []Category
	err := s.Query(ctx, r.DB, r.Map, &categories, query)
	return categories, err
}

func (r *CategoryAdapter) Load(ctx context.Context, id string) (*Category, error) {
	var categories []Category
	query := fmt.Sprintf("select %s from categories where id = %s limit 1", r.Fields, r.BuildParam(1))
	err := s.Query(ctx, r.DB, r.Map, &categories, query, id)
	if err != nil {
		return nil, err
	}
	if len(categories) > 0 {
		return &categories[0], nil
	}
	return nil, nil
}

func (r *CategoryAdapter) Create(ctx context.Context, category *Category) (int64, error) {
	query, args := s.BuildToInsertWithVersion("categories", category, r.VersionIndex, r.BuildParam, true, nil, r.Schema)
	tx := s.GetTx(ctx, r.DB)
	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func (r *CategoryAdapter) Update(ctx context.Context, category *Category) (int64, error) {
	query, args := s.BuildToUpdateWithVersion("categories", category, r.VersionIndex, r.BuildParam, true, nil, r.Schema)
	tx := s.GetTx(ctx, r.DB)
	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return -1, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}
	if rowsAffected <= 0 {
		exist, err := s.Exist(ctx, tx, fmt.Sprintf("select id from categories where id = %s limit 1", r.BuildParam(1)), category.Id)
		if err != nil {
			return -1, err
		}
		if exist {
			return -1, nil
		}
		return 0, nil
	}
	return rowsAffected, nil
}

func (r *CategoryAdapter) Patch(ctx context.Context, category map[string]interface{}) (int64, error) {
	colMap := s.JSONToColumns(category, r.JsonColumnMap)
	query, args := s.BuildToPatchWithVersion("categories", colMap, r.Keys, r.BuildParam, nil, "version", r.Schema.Fields)
	tx := s.GetTx(ctx, r.DB)
	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return -1, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}
	if rowsAffected <= 0 {
		exist, err := s.Exist(ctx, tx, fmt.Sprintf("select id from categories where id = %s limit 1", r.BuildParam(1)), category["id"])
		if err != nil {
			return -1, err
		}
		if exist {
			return -1, nil
		}
		return 0, nil
	}
	return rowsAffected, nil
}

func (r *CategoryAdapter) Delete(ctx context.Context, id string) (int64, error) {
	query := fmt.Sprintf("delete from categories where id = %s", r.BuildParam(1))
	tx := s.GetTx(ctx, r.DB)
	res, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func (r *CategoryAdapter) Search(ctx context.Context, filter *CategoryFilter, limit int64, offset int64) ([]Category, int64, error) {
	var categories []Category
	if limit <= 0 {
		return categories, 0, nil
	}
	query, params := r.BuildQuery(filter)
	pagingQuery := s.BuildPagingQuery(query, limit, offset)
	countQuery := s.BuildCountQuery(query)

	row := r.DB.QueryRowContext(ctx, countQuery, params...)
	if row.Err() != nil {
		return categories, 0, row.Err()
	}
	var total int64
	err := row.Scan(&total)
	if err != nil || total == 0 {
		return categories, total, err
	}

	err = s.Query(ctx, r.DB, r.Map, &categories, pagingQuery, params...)
	return categories, total, err
}

func BuildQuery(filter *CategoryFilter) (string, []interface{}) {
	query := "select * from categories"
	where, params := BuildFilter(filter)
	if len(where) > 0 {
		query = query + " where " + where
	}
	return query, params
}
func BuildFilter(filter *CategoryFilter) (string, []interface{}) {
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
