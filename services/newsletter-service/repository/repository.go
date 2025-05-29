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
	EditorID  string    `json:"editor_id"`
}

// UpdateNewsletterInput represents the input for updating a newsletter.
type UpdateNewsletterInput struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// Post represents a post entity in the database.
type Post struct {
	ID           int       `json:"id"`
	NewsletterID int       `json:"newsletter_id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"created_at"`
	Published    bool      `json:"published"`
}

// Repository defines the methods for interacting with the database.
type Repository interface {
	Save(ctx context.Context, n *Newsletter) error
	FindAll(ctx context.Context) ([]Newsletter, error)
	FindByID(ctx context.Context, id string) (*Newsletter, error) // New method
	Update(ctx context.Context, id string, input UpdateNewsletterInput) (*Newsletter, error)
	Delete(ctx context.Context, id string) error // New method

	// Post-related methods
	CreatePost(ctx context.Context, p *Post) error
	FindPostsByNewsletterID(ctx context.Context, newsletterID string) ([]Post, error)
	UpdatePost(ctx context.Context, id string, p *Post) error
	DeletePost(ctx context.Context, id string) error
	PublishPost(ctx context.Context, postID string) error
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
	query := `INSERT INTO newsletters (id, subject, body, created_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, n.ID, n.Subject, n.Body, n.CreatedAt)
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

// This method retrieves a specific newsletter by its ID:
func (r *postgresRepository) FindByID(ctx context.Context, id string) (*Newsletter, error) {
	query := `SELECT id, subject, body, created_at FROM newsletters WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var n Newsletter
	if err := row.Scan(&n.ID, &n.Subject, &n.Body, &n.CreatedAt); err != nil {
		return nil, err
	}
	return &n, nil
}

// Delete removes a newsletter from the database by its ID.
func (r *postgresRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM newsletters WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *postgresRepository) CreatePost(ctx context.Context, p *Post) error {
	query := `INSERT INTO posts (newsletter_id, title, content, created_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, p.NewsletterID, p.Title, p.Content, p.CreatedAt)
	return err
}

func (r *postgresRepository) FindPostsByNewsletterID(ctx context.Context, newsletterID string) ([]Post, error) {
	query := `SELECT id, newsletter_id, title, content, created_at FROM posts WHERE newsletter_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, newsletterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.NewsletterID, &p.Title, &p.Content, &p.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func (r *postgresRepository) UpdatePost(ctx context.Context, id string, p *Post) error {
	query := `UPDATE posts SET title = $1, content = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, p.Title, p.Content, id)
	return err
}

func (r *postgresRepository) DeletePost(ctx context.Context, id string) error {
	query := `DELETE FROM posts WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *postgresRepository) PublishPost(ctx context.Context, id string) error {
	query := `UPDATE posts SET published = TRUE WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
