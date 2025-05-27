package repository

import (
	"context"
	"database/sql"
	"time"
)

// Newsletter represents a newsletter entity in the database.
type Newsletter struct {
	ID        string    `json:"id"`
	Subject   string    `json:"subject"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

// UpdateNewsletterInput represents the input for updating a newsletter.
type UpdateNewsletterInput struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// Repository defines the methods for interacting with the database.
type Repository interface {
	Save(ctx context.Context, n *Newsletter) error
	FindAll(ctx context.Context) ([]Newsletter, error)
	Update(ctx context.Context, id string, input UpdateNewsletterInput) (*Newsletter, error)
}

type postgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new instance of the repository.
func NewPostgresRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

// Save inserts a new newsletter into the database.
func (r *postgresRepository) Save(ctx context.Context, n *Newsletter) error {
	query := `INSERT INTO newsletters (subject, body, created_at) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, n.Subject, n.Body, n.CreatedAt)
	return err
}

// FindAll retrieves all newsletters from the database.
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

// Update modifies an existing newsletter in the database.
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
