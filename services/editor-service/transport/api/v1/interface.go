package v1

import (
	"context"

	svcmodel "editor-service/service/model"
)

type Service interface {
	CreateEditor(ctx context.Context, Editor svcmodel.Editor) error
	ListEditors(ctx context.Context) []svcmodel.Editor
	GetEditor(ctx context.Context, email string) (svcmodel.Editor, error)
	UpdateEditor(ctx context.Context, email string, Editor svcmodel.Editor) (svcmodel.Editor, error)
	DeleteEditor(ctx context.Context, email string) error
}
