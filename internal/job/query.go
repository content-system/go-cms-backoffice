package job

import "context"

type JobQuery interface {
	Load(ctx context.Context, id string) (*Job, error)
	Search(ctx context.Context, filter *JobFilter, limit int64, offset int64) ([]Job, int64, error)
}
