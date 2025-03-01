package job

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"

	s "github.com/core-go/sql"
)

func NewJobAdapter(db *sql.DB, buildQuery func(*JobFilter) (string, []interface{}), toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}) (*JobAdapter, error) {
	jobType := reflect.TypeOf(Job{})
	params, err := s.CreateParams(jobType, db)
	if err != nil {
		return nil, err
	}
	return &JobAdapter{DB: db, Params: params, BuildQuery: buildQuery, Array: toArray}, nil
}

type JobAdapter struct {
	DB         *sql.DB
	BuildQuery func(*JobFilter) (string, []interface{})
	*s.Params
	Array func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	}
}

func (r *JobAdapter) All(ctx context.Context) ([]Job, error) {
	query := `select * from jobs`
	var jobs []Job
	err := s.Query(ctx, r.DB, r.Map, &jobs, query)
	return jobs, err
}

func (r *JobAdapter) Load(ctx context.Context, id string) (*Job, error) {
	var jobs []Job
	query := fmt.Sprintf("select %s from jobs where id = %s limit 1", r.Fields, r.BuildParam(1))
	err := s.QueryWithArray(ctx, r.DB, r.Map, &jobs, r.Array, query, id)
	if err != nil {
		return nil, err
	}
	if len(jobs) > 0 {
		return &jobs[0], nil
	}
	return nil, nil
}

func (r *JobAdapter) Create(ctx context.Context, job *Job) (int64, error) {
	query, args := s.BuildToInsertWithArray("jobs", job, r.BuildParam, true, r.Array, r.Schema)
	tx := s.GetTx(ctx, r.DB)
	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func (r *JobAdapter) Update(ctx context.Context, job *Job) (int64, error) {
	query, args := s.BuildToUpdateWithArray("jobs", job, r.BuildParam, true, r.Array, r.Schema)
	tx := s.GetTx(ctx, r.DB)
	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func (r *JobAdapter) Patch(ctx context.Context, job map[string]interface{}) (int64, error) {
	colMap := s.JSONToColumns(job, r.JsonColumnMap)
	query, args := s.BuildToPatchWithArray("jobs", colMap, r.Keys, r.BuildParam, r.Array)
	tx := s.GetTx(ctx, r.DB)
	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func (r *JobAdapter) Delete(ctx context.Context, id string) (int64, error) {
	query := fmt.Sprintf("delete from jobs where id = %s", r.BuildParam(1))
	tx := s.GetTx(ctx, r.DB)
	res, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func (r *JobAdapter) Search(ctx context.Context, filter *JobFilter, limit int64, offset int64) ([]Job, int64, error) {
	var jobs []Job
	if limit <= 0 {
		return jobs, 0, nil
	}
	query, params := r.BuildQuery(filter)
	pagingQuery := s.BuildPagingQuery(query, limit, offset)
	countQuery := s.BuildCountQuery(query)

	row := r.DB.QueryRowContext(ctx, countQuery, params...)
	if row.Err() != nil {
		return jobs, 0, row.Err()
	}
	var total int64
	err := row.Scan(&total)
	if err != nil || total == 0 {
		return jobs, total, err
	}

	err = s.QueryWithArray(ctx, r.DB, r.Map, &jobs, r.Array, pagingQuery, params...)
	return jobs, total, err
}

func BuildQuery(filter *JobFilter) (string, []interface{}) {
	query := "select * from jobs"
	where, params := BuildFilter(filter)
	if len(where) > 0 {
		query = query + " where " + where
	}
	return query, params
}
func BuildFilter(filter *JobFilter) (string, []interface{}) {
	buildParam := s.BuildDollarParam
	var where []string
	var params []interface{}
	i := 1
	if len(filter.Id) > 0 {
		params = append(params, filter.Id)
		where = append(where, fmt.Sprintf(`id = %s`, buildParam(i)))
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
