package content

import "context"

type ContentRepository interface {
	Load(ctx context.Context, id string, lang string) (*Content, error)
	Create(ctx context.Context, content *Content) (int64, error)
	Update(ctx context.Context, content *Content) (int64, error)
	Patch(ctx context.Context, content map[string]interface{}) (int64, error)
	Delete(ctx context.Context, id string, lang string) (int64, error)
	Search(ctx context.Context, filter *ContentFilter, limit int64, offset int64) ([]Content, int64, error)
}
