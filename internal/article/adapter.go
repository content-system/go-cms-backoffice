package article

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"

	s "github.com/core-go/sql"
)

func NewArticleAdapter(db *sql.DB, buildQuery func(*ArticleFilter) (string, []interface{}), toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}) (*ArticleAdapter, error) {
	articleType := reflect.TypeOf(Article{})
	params, err := s.CreateParams(articleType, db)
	if err != nil {
		return nil, err
	}
	return &ArticleAdapter{DB: db, Params: params, BuildQuery: buildQuery, Array: toArray}, nil
}

type ArticleAdapter struct {
	DB         *sql.DB
	BuildQuery func(*ArticleFilter) (string, []interface{})
	*s.Params
	Array func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	}
}

func (r *ArticleAdapter) All(ctx context.Context) ([]Article, error) {
	query := `select * from articles`
	var articles []Article
	err := s.Query(ctx, r.DB, r.Map, &articles, query)
	return articles, err
}

func (r *ArticleAdapter) Load(ctx context.Context, id string) (*Article, error) {
	var articles []Article
	query := fmt.Sprintf("select %s from articles where id = %s limit 1", r.Fields, r.BuildParam(1))
	err := s.QueryWithArray(ctx, r.DB, r.Map, &articles, r.Array, query, id)
	if err != nil {
		return nil, err
	}
	if len(articles) > 0 {
		return &articles[0], nil
	}
	return nil, nil
}

func (r *ArticleAdapter) Create(ctx context.Context, article *Article) (int64, error) {
	query, args := s.BuildToInsertWithArray("articles", article, r.BuildParam, true, r.Array, r.Schema)
	tx := s.GetTx(ctx, r.DB)
	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func (r *ArticleAdapter) Update(ctx context.Context, article *Article) (int64, error) {
	query, args := s.BuildToUpdateWithArray("articles", article, r.BuildParam, true, r.Array, r.Schema)
	tx := s.GetTx(ctx, r.DB)
	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func (r *ArticleAdapter) Patch(ctx context.Context, article map[string]interface{}) (int64, error) {
	colMap := s.JSONToColumns(article, r.JsonColumnMap)
	query, args := s.BuildToPatchWithArray("articles", colMap, r.Keys, r.BuildParam, r.Array)
	tx := s.GetTx(ctx, r.DB)
	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func (r *ArticleAdapter) Delete(ctx context.Context, id string) (int64, error) {
	query := fmt.Sprintf("delete from articles where id = %s", r.BuildParam(1))
	tx := s.GetTx(ctx, r.DB)
	res, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func (r *ArticleAdapter) Search(ctx context.Context, filter *ArticleFilter, limit int64, offset int64) ([]Article, int64, error) {
	var articles []Article
	if limit <= 0 {
		return articles, 0, nil
	}
	query, params := r.BuildQuery(filter)
	pagingQuery := s.BuildPagingQuery(query, limit, offset)
	countQuery := s.BuildCountQuery(query)

	row := r.DB.QueryRowContext(ctx, countQuery, params...)
	if row.Err() != nil {
		return articles, 0, row.Err()
	}
	var total int64
	err := row.Scan(&total)
	if err != nil || total == 0 {
		return articles, total, err
	}

	err = s.QueryWithArray(ctx, r.DB, r.Map, &articles, r.Array, pagingQuery, params...)
	return articles, total, err
}

func BuildQuery(filter *ArticleFilter) (string, []interface{}) {
	query := "select * from articles"
	where, params := BuildFilter(filter)
	if len(where) > 0 {
		query = query + " where " + where
	}
	return query, params
}
func BuildFilter(filter *ArticleFilter) (string, []interface{}) {
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
