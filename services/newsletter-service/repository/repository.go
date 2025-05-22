package newsletter

import (
	"context"
	"database/sql"
)

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Save(ctx context.Context, n *Newsletter) error {
	query := `INSERT INTO newsletters (id, subject, body, created_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, n.ID, n.Subject, n.Body, n.CreatedAt)
	return err
}

func (r *postgresRepository) FindAll(ctx context.Context) ([]Newsletter, error) {
	query := `SELECT id, subject, body, created_at FROM newsletters ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var newsletters []Newsletter
	for rows.Next() {
		var n Newsletter
		if err := rows.Scan(&n.ID, &n.Subject, &n.Body, &n.CreatedAt); err != nil {
			return nil, err
		}
		newsletters = append(newsletters, n)
	}
	return newsletters, nil
}

func (r *postgresRepository) Update(ctx context.Context, id string, input UpdateNewsletterInput) (*Newsletter, error) {
	query := `UPDATE newsletters SET subject = $1, body = $2 WHERE id = $3 RETURNING id, subject, body, created_at`
	var n Newsletter
	err := r.db.QueryRowContext(ctx, query, input.Subject, input.Body, id).
		Scan(&n.ID, &n.Subject, &n.Body, &n.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &n, nil
}
