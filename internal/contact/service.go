package contact

import (
	"context"
	"database/sql"

	"github.com/core-go/core/tx"
)

type ContactService interface {
	Load(ctx context.Context, id string) (*Contact, error)
	Create(ctx context.Context, contact *Contact) (int64, error)
	Update(ctx context.Context, contact *Contact) (int64, error)
	Patch(ctx context.Context, contact map[string]interface{}) (int64, error)
	Delete(ctx context.Context, id string) (int64, error)
	Search(ctx context.Context, filter *ContactFilter, limit int64, offset int64) ([]Contact, int64, error)
}

func NewContactService(db *sql.DB, repository ContactRepository) *ContactUseCase {
	return &ContactUseCase{db: db, repository: repository}
}

type ContactUseCase struct {
	db         *sql.DB
	repository ContactRepository
}

func (s *ContactUseCase) Load(ctx context.Context, id string) (*Contact, error) {
	return s.repository.Load(ctx, id)
}
func (s *ContactUseCase) Create(ctx context.Context, contact *Contact) (int64, error) {
	return tx.Execute(ctx, s.db, func(ctx context.Context) (int64, error) {
		return s.repository.Create(ctx, contact)
	})
}
func (s *ContactUseCase) Update(ctx context.Context, contact *Contact) (int64, error) {
	return tx.Execute(ctx, s.db, func(ctx context.Context) (int64, error) {
		return s.repository.Update(ctx, contact)
	})
}
func (s *ContactUseCase) Patch(ctx context.Context, contact map[string]interface{}) (int64, error) {
	return tx.Execute(ctx, s.db, func(ctx context.Context) (int64, error) {
		return s.repository.Patch(ctx, contact)
	})
}
func (s *ContactUseCase) Delete(ctx context.Context, id string) (int64, error) {
	return tx.Execute(ctx, s.db, func(ctx context.Context) (int64, error) {
		return s.repository.Delete(ctx, id)
	})
}
func (s *ContactUseCase) Search(ctx context.Context, filter *ContactFilter, limit int64, offset int64) ([]Contact, int64, error) {
	return s.repository.Search(ctx, filter, limit, offset)
}
