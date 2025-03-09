package category

import "context"

type CategoryRepository interface {
	Load(ctx context.Context, id string) (*Category, error)
	Create(ctx context.Context, category *Category) (int64, error)
	Update(ctx context.Context, category *Category) (int64, error)
	Patch(ctx context.Context, category map[string]interface{}) (int64, error)
	Delete(ctx context.Context, id string) (int64, error)
	Search(ctx context.Context, filter *CategoryFilter, limit int64, offset int64) ([]Category, int64, error)
}
