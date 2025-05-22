package repository

import (
	"context"

	"editor-service/pkg/id"
	dbmodel "editor-service/repository/sql/model"
	"editor-service/repository/sql/query"
	"editor-service/service/model"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EditorRepository struct {
	pool *pgxpool.Pool
}

func NewEditorRepository(pool *pgxpool.Pool) *EditorRepository {
	return &EditorRepository{
		pool: pool,
	}
}

func (r *EditorRepository) ReadEditor(ctx context.Context, EditorID id.Editor) (*model.Editor, error) {
	var Editor dbmodel.Editor
	if err := pgxscan.Get(
		ctx,
		r.pool,
		&Editor,
		query.ReadEditor,
		pgx.NamedArgs{
			"id": EditorID,
		},
	); err != nil {
		return nil, err
	}
	return &model.Editor{}, nil
}

func (r *EditorRepository) ListEditor(ctx context.Context) ([]model.Editor, error) {
	var Editors []dbmodel.Editor
	if err := pgxscan.Select(
		ctx,
		r.pool,
		&Editors,
		query.ListEditor,
	); err != nil {
		return nil, err
	}
	response := make([]model.Editor, len(Editors))
	for i, Editor := range Editors {
		response[i] = model.Editor{
			ID: Editor.ID,
		}
	}
	return response, nil
}
