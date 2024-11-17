package contact

import (
	"context"
	"database/sql"
	"strings"

	s "github.com/core-go/sql"
)

func NewContactAdapter(db *sql.DB) *ContactAdapter {
	buildParam := s.GetBuild(db)
	return &ContactAdapter{DB: db, BuildParam: buildParam}
}

type ContactAdapter struct {
	DB         *sql.DB
	BuildParam func(int) string
}

func (r *ContactAdapter) Submit(ctx context.Context, contact *Contact) (int64, error) {
	query := `
		insert into contacts (
			id,
			name,
			country,
			company,
			job_title,
			email,
			phone,
			submitted_at,
			message)
		values (
			$1,
			$2,
			$3, 
			$4,
			$5,
			$6,
			$7,
			$8,
			$9)`
	stmt, err := r.DB.Prepare(query)
	if err != nil {
		return -1, nil
	}
	res, err := stmt.ExecContext(ctx,
		contact.Id,
		contact.Name,
		contact.Country,
		contact.Company,
		contact.JobTitle,
		contact.Email,
		contact.Phone,
		contact.SubmittedAt,
		contact.Message)
	if err != nil {
		if strings.Index(err.Error(), "duplicate key") >= 0 {
			return -1, nil
		}
		return -1, err
	}
	return res.RowsAffected()
}
