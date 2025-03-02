package content

import (
	"context"
	"database/sql"

	"github.com/core-go/core/tx"
)

type ContentService interface {
	Load(ctx context.Context, id string, lang string) (*Content, error)
	Create(ctx context.Context, content *Content) (int64, error)
	Update(ctx context.Context, content *Content) (int64, error)
	Patch(ctx context.Context, content map[string]interface{}) (int64, error)
	Delete(ctx context.Context, id string, lang string) (int64, error)
	Search(ctx context.Context, filter *ContentFilter, limit int64, offset int64) ([]Content, int64, error)
}

func NewContentService(repository ContentRepository) *ContentUseCase {
	return &ContentUseCase{repository: repository}
}

type ContentUseCase struct {
	db         *sql.DB
	repository ContentRepository
}

func (s *ContentUseCase) Load(ctx context.Context, id string, lang string) (*Content, error) {
	return s.repository.Load(ctx, id, lang)
}
func (s *ContentUseCase) Create(ctx context.Context, content *Content) (int64, error) {
	return tx.Execute(ctx, s.db, func(ctx context.Context) (int64, error) {
		return s.repository.Create(ctx, content)
	})
}
func (s *ContentUseCase) Update(ctx context.Context, content *Content) (int64, error) {
	return tx.Execute(ctx, s.db, func(ctx context.Context) (int64, error) {
		return s.repository.Update(ctx, content)
	})
}
func (s *ContentUseCase) Patch(ctx context.Context, content map[string]interface{}) (int64, error) {
	return tx.Execute(ctx, s.db, func(ctx context.Context) (int64, error) {
		return s.repository.Patch(ctx, content)
	})
}
func (s *ContentUseCase) Delete(ctx context.Context, id string, lang string) (int64, error) {
	return tx.Execute(ctx, s.db, func(ctx context.Context) (int64, error) {
		return s.repository.Delete(ctx, id, lang)
	})
}
func (s *ContentUseCase) Search(ctx context.Context, filter *ContentFilter, limit int64, offset int64) ([]Content, int64, error) {
	return s.repository.Search(ctx, filter, limit, offset)
}
