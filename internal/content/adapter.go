package content

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"

	s "github.com/core-go/sql"
)

func NewContentAdapter(db *sql.DB, buildQuery func(*ContentFilter) (string, []interface{}), toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}) (*ContentAdapter, error) {
	contentType := reflect.TypeOf(Content{})
	params, err := s.CreateParams(contentType, db)
	if err != nil {
		return nil, err
	}
	return &ContentAdapter{DB: db, Params: params, BuildQuery: buildQuery, Array: toArray}, nil
}

type ContentAdapter struct {
	DB         *sql.DB
	BuildQuery func(*ContentFilter) (string, []interface{})
	*s.Params
	Array func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	}
}

func (r *ContentAdapter) All(ctx context.Context) ([]Content, error) {
	query := `select * from contents`
	var contents []Content
	err := s.Query(ctx, r.DB, r.Map, &contents, query)
	return contents, err
}

func (r *ContentAdapter) Load(ctx context.Context, id string, lang string) (*Content, error) {
	var contents []Content
	query := fmt.Sprintf("select %s from contents where id = %s and lang = %s limit 1", r.Fields, r.BuildParam(1), r.BuildParam(2))
	err := s.QueryWithArray(ctx, r.DB, r.Map, &contents, r.Array, query, id, lang)
	if err != nil {
		return nil, err
	}
	if len(contents) > 0 {
		return &contents[0], nil
	}
	return nil, nil
}

func (r *ContentAdapter) Create(ctx context.Context, content *Content) (int64, error) {
	query, args := s.BuildToInsertWithArray("contents", content, r.BuildParam, true, r.Array, r.Schema)
	tx := s.GetTx(ctx, r.DB)
	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func (r *ContentAdapter) Update(ctx context.Context, content *Content) (int64, error) {
	query, args := s.BuildToUpdateWithArray("contents", content, r.BuildParam, true, r.Array, r.Schema)
	tx := s.GetTx(ctx, r.DB)
	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func (r *ContentAdapter) Patch(ctx context.Context, content map[string]interface{}) (int64, error) {
	colMap := s.JSONToColumns(content, r.JsonColumnMap)
	query, args := s.BuildToPatchWithArray("contents", colMap, r.Keys, r.BuildParam, r.Array)
	tx := s.GetTx(ctx, r.DB)
	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func (r *ContentAdapter) Delete(ctx context.Context, id string, lang string) (int64, error) {
	query := fmt.Sprintf("delete from contents where id = %s and lang = %s", r.BuildParam(1), r.BuildParam(2))
	tx := s.GetTx(ctx, r.DB)
	res, err := tx.ExecContext(ctx, query, id, lang)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func (r *ContentAdapter) Search(ctx context.Context, filter *ContentFilter, limit int64, offset int64) ([]Content, int64, error) {
	var contents []Content
	if limit <= 0 {
		return contents, 0, nil
	}
	query, params := r.BuildQuery(filter)
	pagingQuery := s.BuildPagingQuery(query, limit, offset)
	countQuery := s.BuildCountQuery(query)

	row := r.DB.QueryRowContext(ctx, countQuery, params...)
	if row.Err() != nil {
		return contents, 0, row.Err()
	}
	var total int64
	err := row.Scan(&total)
	if err != nil || total == 0 {
		return contents, total, err
	}

	err = s.QueryWithArray(ctx, r.DB, r.Map, &contents, r.Array, pagingQuery, params...)
	return contents, total, err
}

func BuildQuery(filter *ContentFilter) (string, []interface{}) {
	query := "select * from contents"
	where, params := BuildFilter(filter)
	if len(where) > 0 {
		query = query + " where " + where
	}
	return query, params
}
func BuildFilter(filter *ContentFilter) (string, []interface{}) {
	buildParam := s.BuildDollarParam
	var where []string
	var params []interface{}
	i := 1
	if len(filter.Id) > 0 {
		params = append(params, filter.Id)
		where = append(where, fmt.Sprintf(`id = %s`, buildParam(i)))
		i++
	}
	if len(filter.Lang) > 0 {
		params = append(params, filter.Id)
		where = append(where, fmt.Sprintf(`lang = %s`, buildParam(i)))
		i++
	}
	if filter.PublishedAt != nil {
		if filter.PublishedAt.Min != nil {
			params = append(params, filter.PublishedAt.Min)
			where = append(where, fmt.Sprintf(`published_at >= %s`, buildParam(i)))
			i++
		}
		if filter.PublishedAt.Max != nil {
			params = append(params, filter.PublishedAt.Max)
			where = append(where, fmt.Sprintf(`published_at <= %s`, buildParam(i)))
			i++
		}
	}
	if len(filter.Title) > 0 {
		q := filter.Title + "%"
		params = append(params, q)
		where = append(where, fmt.Sprintf(`title like %s`, buildParam(i)))
		i++
	}
	if len(where) > 0 {
		return strings.Join(where, " and "), params
	}
	return "", params
}
