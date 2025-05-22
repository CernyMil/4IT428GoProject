package repository

import (
	"context"
	"editor-service/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EditorRepository interface {
	CreateEditor(ctx context.Context, editor *models.Editor) error
	GetEditorByEmail(ctx context.Context, email string) (*models.Editor, error)
}

type editorRepository struct {
	db *pgxpool.Pool
}

func NewEditorRepository(db *pgxpool.Pool) EditorRepository {
	return &editorRepository{db: db}
}

func (r *editorRepository) CreateEditor(ctx context.Context, editor *models.Editor) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO editors (id, email, first_name, last_name, created_at) 
         VALUES ($1, $2, $3, $4, $5)`,
		editor.ID, editor.Email, editor.FirstName, editor.LastName, editor.CreatedAt,
	)
	return err
}

func (r *editorRepository) GetEditorByEmail(ctx context.Context, email string) (*models.Editor, error) {
	row := r.db.QueryRow(ctx, `SELECT id, email, first_name, last_name, created_at FROM editors WHERE email=$1`, email)

	var editor models.Editor
	err := row.Scan(&editor.ID, &editor.Email, &editor.FirstName, &editor.LastName, &editor.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &editor, nil
}
