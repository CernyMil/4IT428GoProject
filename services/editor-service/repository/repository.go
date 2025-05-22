package repository

import "github.com/jackc/pgx/v5/pgxpool"

type Repository struct {
	*EditorRepository
}

func New(pool *pgxpool.Pool) (*Repository, error) {
	return &Repository{
		EditorRepository: NewEditorRepository(pool),
	}, nil
}
