package job

import "context"

type JobRepository interface {
	Load(ctx context.Context, id string) (*Job, error)
	Create(ctx context.Context, job *Job) (int64, error)
	Update(ctx context.Context, job *Job) (int64, error)
	Patch(ctx context.Context, job map[string]interface{}) (int64, error)
	Delete(ctx context.Context, id string) (int64, error)
}
