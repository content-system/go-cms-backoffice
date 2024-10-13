package job

import (
	"context"
	"database/sql"
	"github.com/core-go/core/tx"
)

type JobService interface {
	Load(ctx context.Context, id string) (*Job, error)
	Create(ctx context.Context, job *Job) (int64, error)
	Update(ctx context.Context, job *Job) (int64, error)
	Patch(ctx context.Context, job map[string]interface{}) (int64, error)
	Delete(ctx context.Context, id string) (int64, error)
}

func NewJobService(repository JobRepository) *JobUseCase {
	return &JobUseCase{repository: repository}
}

type JobUseCase struct {
	db         *sql.DB
	repository JobRepository
}

func (s *JobUseCase) Load(ctx context.Context, id string) (*Job, error) {
	return s.repository.Load(ctx, id)
}
func (s *JobUseCase) Create(ctx context.Context, job *Job) (int64, error) {
	return tx.Execute(ctx, s.db, func(ctx context.Context) (int64, error) {
		return s.repository.Create(ctx, job)
	})
}
func (s *JobUseCase) Update(ctx context.Context, job *Job) (int64, error) {
	return tx.Execute(ctx, s.db, func(ctx context.Context) (int64, error) {
		return s.repository.Update(ctx, job)
	})
}
func (s *JobUseCase) Patch(ctx context.Context, job map[string]interface{}) (int64, error) {
	return tx.Execute(ctx, s.db, func(ctx context.Context) (int64, error) {
		return s.repository.Patch(ctx, job)
	})
}
func (s *JobUseCase) Delete(ctx context.Context, id string) (int64, error) {
	return tx.Execute(ctx, s.db, func(ctx context.Context) (int64, error) {
		return s.repository.Delete(ctx, id)
	})
}
