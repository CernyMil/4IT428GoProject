package service

import (
	"context"
	"editor-service/models"
	uuidutil "editor-service/pkg/id"
	"editor-service/repository"
	transport "editor-service/transport/middleware"
	"fmt"
	"time"

	"firebase.google.com/go/auth"
)

type EditorService struct {
	repo repository.EditorRepositoryInterface
	auth *transport.FirebaseAuth
}

func NewEditorService(repo repository.EditorRepositoryInterface, auth *transport.FirebaseAuth) *EditorService {
	return &EditorService{repo: repo, auth: auth}
}

func (s *EditorService) SignUp(ctx context.Context, email, password, firstName, lastName string) error {
	params := (&auth.UserToCreate{}).
		Email(email).
		Password(password)
	userRecord, err := s.auth.Client.CreateUser(ctx, params)
	if err != nil {
		return fmt.Errorf("firebase create user error: %w", err)
	}

	editor := &models.Editor{
		ID:        uuidutil.NewUUID(),
		Email:     userRecord.Email,
		FirstName: firstName,
		LastName:  lastName,
		CreatedAt: time.Now(),
	}
	return s.repo.CreateEditor(ctx, editor)
}
