package repository

import (
	"context"
	"editor-service/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EditorRepositoryInterface interface {
	CreateEditor(ctx context.Context, editor *models.Editor) error
	GetEditorByEmail(ctx context.Context, email string) (*models.Editor, error)
}

type EditorRepository struct {
	db *pgxpool.Pool
}

func NewPgxEditorRepository(db *pgxpool.Pool) EditorRepositoryInterface {
	return &EditorRepository{db: db}
}

func (r *EditorRepository) CreateEditor(ctx context.Context, editor *models.Editor) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO editors (id, email, first_name, last_name, created_at, hashed_password) 
         VALUES ($1, $2, $3, $4, $5, $6)`,
		editor.ID, editor.Email, editor.FirstName, editor.LastName, editor.CreatedAt, editor.HashedPassword,
	)
	return err
}

func (r *EditorRepository) GetEditorByEmail(ctx context.Context, email string) (*models.Editor, error) {
	row := r.db.QueryRow(ctx, `SELECT id, email, first_name, last_name, created_at, hashed_password FROM editors WHERE email=$1`, email)

	var editor models.Editor
	err := row.Scan(&editor.ID, &editor.Email, &editor.FirstName, &editor.LastName, &editor.CreatedAt, &editor.HashedPassword)
	if err != nil {
		return nil, err
	}

	return &editor, nil
}
