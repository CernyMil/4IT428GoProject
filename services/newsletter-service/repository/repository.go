package repository

import (
	"context"
	"newsletter-service/pkg/id"
	svcmodel "newsletter-service/service/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresRepository creates a new instance of postgresRepository.
func NewPostgresRepository(pool *pgxpool.Pool) (*postgresRepository, error) {
	return &postgresRepository{
		pool: pool,
	}, nil
}

// Save inserts a new newsletter into the database.
func (r *postgresRepository) Save(ctx context.Context, n *svcmodel.Newsletter) error {
	query := `INSERT INTO newsletters (id, subject, body, created_at) VALUES ($1, $2, $3, $4)`
	_, err := r.pool.Exec(ctx, query, n.ID, n.Title, n.Description, n.CreatedAt)
	return err
}

// FindAll retrieves all newsletters from the database.
func (r *postgresRepository) FindAll(ctx context.Context) ([]svcmodel.Newsletter, error) {
	query := `SELECT id, subject, body, created_at FROM newsletters ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var newsletters []svcmodel.Newsletter
	for rows.Next() {
		var n svcmodel.Newsletter
		var idStr string // Temporary variable to hold the UUID as a string
		if err := rows.Scan(&idStr, &n.Title, &n.Description, &n.CreatedAt); err != nil {
			return nil, err
		}
		// Parse the UUID string into id.Newsletter
		parsedID, err := uuid.Parse(idStr)
		if err != nil {
			return nil, err
		}
		n.ID = id.Newsletter(parsedID) // Convert uuid.UUID to id.Newsletter
		newsletters = append(newsletters, n)
	}
	return newsletters, nil
}

// Update modifies an existing newsletter in the database.
func (r *postgresRepository) Update(ctx context.Context, id id.Newsletter, input svcmodel.UpdateNewsletterInput) (*svcmodel.Newsletter, error) {
	query := `UPDATE newsletters SET subject = $1, body = $2 WHERE id = $3 RETURNING id, subject, body, created_at`
	var n svcmodel.Newsletter
	err := r.pool.QueryRow(ctx, query, input.Title, input.Description, id).
		Scan(&n.ID, &n.Title, &n.Description, &n.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &n, nil
}

// FindByID retrieves a specific newsletter by its ID.
func (r *postgresRepository) FindByID(ctx context.Context, id string) (*svcmodel.Newsletter, error) {
	query := `SELECT id, subject, body, created_at FROM newsletters WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, id)

	var n svcmodel.Newsletter
	if err := row.Scan(&n.ID, &n.Title, &n.Description, &n.CreatedAt); err != nil {
		return nil, err
	}
	return &n, nil
}

// Delete removes a newsletter from the database by its ID.
func (r *postgresRepository) Delete(ctx context.Context, id id.Newsletter) error {
	query := `DELETE FROM newsletters WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

// CreatePost inserts a new post into the database.
func (r *postgresRepository) CreatePost(ctx context.Context, p *svcmodel.Post) error {
	query := `INSERT INTO posts (newsletter_id, title, content, created_at) VALUES ($1, $2, $3, $4)`
	_, err := r.pool.Exec(ctx, query, p.NewsletterID, p.Title, p.Content, p.CreatedAt)
	return err
}

// FindPostsByNewsletterID retrieves all posts for a specific newsletter.
func (r *postgresRepository) FindPostsByNewsletterID(ctx context.Context, newsletterID id.Newsletter) ([]svcmodel.Post, error) {
	query := `SELECT id, newsletter_id, title, content, created_at, published FROM posts WHERE newsletter_id = $1 ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, query, newsletterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []svcmodel.Post
	for rows.Next() {
		var p svcmodel.Post
		if err := rows.Scan(&p.ID, &p.NewsletterID, &p.Title, &p.Content, &p.CreatedAt, &p.Published); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

// UpdatePost modifies an existing post in the database.
func (r *postgresRepository) UpdatePost(ctx context.Context, id id.Post, p *svcmodel.Post) error {
	query := `UPDATE posts SET title = $1, content = $2 WHERE id = $3`
	_, err := r.pool.Exec(ctx, query, p.Title, p.Content, id)
	return err
}

// DeletePost removes a post from the database by its ID.
func (r *postgresRepository) DeletePost(ctx context.Context, id id.Post) error {
	query := `DELETE FROM posts WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

// PublishPost marks a post as published in the database.
func (r *postgresRepository) PublishPost(ctx context.Context, id id.Post) error {
	query := `UPDATE posts SET published = TRUE WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}
