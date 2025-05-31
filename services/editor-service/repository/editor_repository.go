package repository

import (
	"context"
	"editor-service/models"
	"editor-service/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EditorRepository struct {
	db *pgxpool.Pool
}

func NewPgxEditorRepository(db *pgxpool.Pool) service.EditorRepositoryInterface {
	return &EditorRepository{db: db}
}

func (r *EditorRepository) CreateEditor(ctx context.Context, editor *models.Editor) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO editors (id, email, first_name, last_name, created_at) 
         VALUES ($1, $2, $3, $4, $5)`,
		editor.ID, editor.Email, editor.FirstName, editor.LastName, editor.CreatedAt,
	)
	return err
}

func (r *EditorRepository) GetEditorByEmail(ctx context.Context, email string) (*models.Editor, error) {
	row := r.db.QueryRow(ctx, `SELECT id, email, first_name, last_name, created_at FROM editors WHERE email=$1`, email)

	var editor models.Editor
	err := row.Scan(&editor.ID, &editor.Email, &editor.FirstName, &editor.LastName, &editor.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &editor, nil
}
