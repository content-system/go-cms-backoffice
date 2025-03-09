package category

import (
	"context"
	"database/sql"

	"github.com/core-go/core/tx"
)

type CategoryService interface {
	Load(ctx context.Context, id string) (*Category, error)
	Create(ctx context.Context, category *Category) (int64, error)
	Update(ctx context.Context, category *Category) (int64, error)
	Patch(ctx context.Context, category map[string]interface{}) (int64, error)
	Delete(ctx context.Context, id string) (int64, error)
	Search(ctx context.Context, filter *CategoryFilter, limit int64, offset int64) ([]Category, int64, error)
}

func NewCategoryService(db *sql.DB, repository CategoryRepository) *CategoryUseCase {
	return &CategoryUseCase{db: db, repository: repository}
}

type CategoryUseCase struct {
	db         *sql.DB
	repository CategoryRepository
}

func (s *CategoryUseCase) Load(ctx context.Context, id string) (*Category, error) {
	return s.repository.Load(ctx, id)
}
func (s *CategoryUseCase) Create(ctx context.Context, category *Category) (int64, error) {
	return tx.Execute(ctx, s.db, func(ctx context.Context) (int64, error) {
		return s.repository.Create(ctx, category)
	})
}
func (s *CategoryUseCase) Update(ctx context.Context, category *Category) (int64, error) {
	return tx.Execute(ctx, s.db, func(ctx context.Context) (int64, error) {
		return s.repository.Update(ctx, category)
	})
}
func (s *CategoryUseCase) Patch(ctx context.Context, category map[string]interface{}) (int64, error) {
	return tx.Execute(ctx, s.db, func(ctx context.Context) (int64, error) {
		return s.repository.Patch(ctx, category)
	})
}
func (s *CategoryUseCase) Delete(ctx context.Context, id string) (int64, error) {
	return tx.Execute(ctx, s.db, func(ctx context.Context) (int64, error) {
		return s.repository.Delete(ctx, id)
	})
}
func (s *CategoryUseCase) Search(ctx context.Context, filter *CategoryFilter, limit int64, offset int64) ([]Category, int64, error) {
	return s.repository.Search(ctx, filter, limit, offset)
}
