package article

import (
	"context"
	"database/sql"

	"github.com/core-go/core/tx"
)

type ArticleService interface {
	Load(ctx context.Context, id string) (*Article, error)
	Create(ctx context.Context, article *Article) (int64, error)
	Update(ctx context.Context, article *Article) (int64, error)
	Patch(ctx context.Context, article map[string]interface{}) (int64, error)
	Delete(ctx context.Context, id string) (int64, error)
	Search(ctx context.Context, filter *ArticleFilter, limit int64, offset int64) ([]Article, int64, error)
}

func NewArticleService(repository ArticleRepository) *ArticleUseCase {
	return &ArticleUseCase{repository: repository}
}

type ArticleUseCase struct {
	db         *sql.DB
	repository ArticleRepository
}

func (s *ArticleUseCase) Load(ctx context.Context, id string) (*Article, error) {
	return s.repository.Load(ctx, id)
}
func (s *ArticleUseCase) Create(ctx context.Context, article *Article) (int64, error) {
	return tx.Execute(ctx, s.db, func(ctx context.Context) (int64, error) {
		return s.repository.Create(ctx, article)
	})
}
func (s *ArticleUseCase) Update(ctx context.Context, article *Article) (int64, error) {
	return tx.Execute(ctx, s.db, func(ctx context.Context) (int64, error) {
		return s.repository.Update(ctx, article)
	})
}
func (s *ArticleUseCase) Patch(ctx context.Context, article map[string]interface{}) (int64, error) {
	return tx.Execute(ctx, s.db, func(ctx context.Context) (int64, error) {
		return s.repository.Patch(ctx, article)
	})
}
func (s *ArticleUseCase) Delete(ctx context.Context, id string) (int64, error) {
	return tx.Execute(ctx, s.db, func(ctx context.Context) (int64, error) {
		return s.repository.Delete(ctx, id)
	})
}
func (s *ArticleUseCase) Search(ctx context.Context, filter *ArticleFilter, limit int64, offset int64) ([]Article, int64, error) {
	return s.repository.Search(ctx, filter, limit, offset)
}
