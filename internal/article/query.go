package article

import "context"

type ArticleQuery interface {
	Load(ctx context.Context, id string) (*Article, error)
	Search(ctx context.Context, filter *ArticleFilter, limit int64, offset int64) ([]Article, int64, error)
}
