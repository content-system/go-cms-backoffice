package article

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"

	q "github.com/core-go/sql"
)

func NewDraftArticleAdapter(db *sql.DB, buildQuery func(*ArticleFilter) (string, []interface{}), toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}) (*DraftArticleAdapter, error) {
	parameters, err := q.CreateParameters(reflect.TypeOf(Article{}), db)
	if err != nil {
		return nil, err
	}
	return &DraftArticleAdapter{DB: db, Parameters: parameters, BuildQuery: buildQuery, Array: toArray}, nil
}

type DraftArticleAdapter struct {
	DB         *sql.DB
	BuildQuery func(*ArticleFilter) (string, []interface{})
	*q.Parameters
	Array func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	}
}

func (r *DraftArticleAdapter) All(ctx context.Context) ([]Article, error) {
	query := `select * from draft_articles`
	var articles []Article
	err := q.Query(ctx, r.DB, r.Map, &articles, query)
	return articles, err
}

func (r *DraftArticleAdapter) Load(ctx context.Context, id string) (*Article, error) {
	var articles []Article
	query := fmt.Sprintf("select %s from draft_articles where id = %s limit 1", r.Fields, r.BuildParam(1))
	fmt.Println(query)
	err := q.QueryWithArray(ctx, r.DB, r.Map, &articles, r.Array, query, id)
	if err != nil {
		return nil, err
	}
	if len(articles) > 0 {
		return &articles[0], nil
	}
	return nil, nil
}

func (r *DraftArticleAdapter) Create(ctx context.Context, article *Article) (int64, error) {
	query, args := q.BuildToInsertWithArray("draft_articles", article, r.BuildParam, true, r.Array, r.Schema)
	tx := q.GetTx(ctx, r.DB)
	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func (r *DraftArticleAdapter) Update(ctx context.Context, article *Article) (int64, error) {
	query, args := q.BuildToUpdateWithArray("draft_articles", article, r.BuildParam, true, r.Array, r.Schema)
	tx := q.GetTx(ctx, r.DB)
	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func (r *DraftArticleAdapter) Patch(ctx context.Context, article map[string]interface{}) (int64, error) {
	colMap := q.JSONToColumns(article, r.JsonColumnMap)
	query, args := q.BuildToPatchWithArray("draft_articles", colMap, r.Keys, r.BuildParam, r.Array)
	tx := q.GetTx(ctx, r.DB)
	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func (r *DraftArticleAdapter) Delete(ctx context.Context, id string) (int64, error) {
	query := fmt.Sprintf("delete from draft_articles where id = %s", r.BuildParam(1))
	tx := q.GetTx(ctx, r.DB)
	res, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func (r *DraftArticleAdapter) Search(ctx context.Context, filter *ArticleFilter, limit int64, offset int64) ([]Article, int64, error) {
	var articles []Article
	if limit <= 0 {
		return articles, 0, nil
	}
	query, params := r.BuildQuery(filter)
	pagingQuery := q.BuildPagingQuery(query, limit, offset)
	countQuery := q.BuildCountQuery(query)

	row := r.DB.QueryRowContext(ctx, countQuery, params...)
	if row.Err() != nil {
		return articles, 0, row.Err()
	}
	var total int64
	err := row.Scan(&total)
	if err != nil || total == 0 {
		return articles, total, err
	}

	err = q.QueryWithArray(ctx, r.DB, r.Map, &articles, r.Array, pagingQuery, params...)
	return articles, total, err
}

func BuildQuery(filter *ArticleFilter) (string, []interface{}) {
	query := "select * from draft_articles"
	where, params := BuildFilter(filter)
	if len(where) > 0 {
		query = query + " where " + where
	}
	return query, params
}
func BuildFilter(filter *ArticleFilter) (string, []interface{}) {
	buildParam := q.BuildDollarParam
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
