package contact

import "context"

type ContactRepository interface {
	Load(ctx context.Context, id string) (*Contact, error)
	Create(ctx context.Context, contact *Contact) (int64, error)
	Update(ctx context.Context, contact *Contact) (int64, error)
	Patch(ctx context.Context, contact map[string]interface{}) (int64, error)
	Delete(ctx context.Context, id string) (int64, error)
	Search(ctx context.Context, filter *ContactFilter, limit int64, offset int64) ([]Contact, int64, error)
}
