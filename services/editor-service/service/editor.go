package service

import (
	"context"
	"editor-service/models"
	"editor-service/pkg/id"
	"editor-service/repository"
	"editor-service/transport"
	"time"
)

type EditorService struct {
	repo repository.EditorRepository
	auth *transport.FirebaseAuth
}

func NewEditorService(repo repository.EditorRepository, auth *transport.FirebaseAuth) *EditorService {
	return &EditorService{repo: repo, auth: auth}
}

func (s *EditorService) SignUp(ctx context.Context, idToken, firstName, lastName string) error {
	token, err := s.auth.VerifyIDToken(ctx, idToken)
	if err != nil {
		return err
	}

	email := token.Claims["email"].(string)

	editor := &models.Editor{
		ID:        uuidutil.NewUUID(),
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		CreatedAt: time.Now(),
	}

	return s.repo.CreateEditor(ctx, editor)
}
