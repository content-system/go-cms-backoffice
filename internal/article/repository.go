package article

import "context"

type ArticleRepository interface {
	Load(ctx context.Context, id string) (*Article, error)
	Create(ctx context.Context, article *Article) (int64, error)
	Update(ctx context.Context, article *Article) (int64, error)
	Patch(ctx context.Context, article map[string]interface{}) (int64, error)
	Delete(ctx context.Context, id string) (int64, error)
}
