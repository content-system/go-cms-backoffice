package article

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"

	q "github.com/core-go/sql"
)

func NewArticleAdapter(db *sql.DB, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}) (*ArticleAdapter, error) {
	parameters, err := q.CreateParameters(reflect.TypeOf(Article{}), db)
	if err != nil {
		return nil, err
	}
	return &ArticleAdapter{DB: db, Parameters: parameters, Array: toArray}, nil
}

type ArticleAdapter struct {
	DB *sql.DB
	*q.Parameters
	Array func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	}
}

func (r *ArticleAdapter) Exist(ctx context.Context, id string) (bool, error) {
	var exists bool
	query := `select exists (select 1 from articles where id = $1)`
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *ArticleAdapter) Load(ctx context.Context, id string) (*Article, error) {
	var articles []Article
	query := fmt.Sprintf("select %s from articles where id = %s limit 1", r.Fields, r.BuildParam(1))
	err := q.QueryWithArray(ctx, r.DB, r.Map, &articles, r.Array, query, id)
	if err != nil {
		return nil, err
	}
	if len(articles) > 0 {
		return &articles[0], nil
	}
	return nil, nil
}

func (r *ArticleAdapter) Save(ctx context.Context, article *Article) (int64, error) {
	query, args, err := q.BuildToSaveWithArray("articles", article, q.DriverPostgres, r.Array, r.Schema)
	if err != nil {
		return -1, err
	}
	tx := q.GetTx(ctx, r.DB)
	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}
