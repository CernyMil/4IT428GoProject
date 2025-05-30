package service

import (
	"context"
	"editor-service/models"
	uuidutil "editor-service/pkg/id"
	transport "editor-service/transport/middleware"
	"fmt"
	"time"

	"firebase.google.com/go/auth"
)

type EditorRepositoryInterface interface {
	CreateEditor(ctx context.Context, editor *models.Editor) error
	GetEditorByEmail(ctx context.Context, email string) (*models.Editor, error)
}

type EditorService struct {
	repo EditorRepositoryInterface
	auth *transport.FirebaseAuth
}

func NewEditorService(repo EditorRepositoryInterface, auth *transport.FirebaseAuth) *EditorService {
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

func (s *EditorService) GetByEmail(ctx context.Context, email string) (*models.Editor, error) {
	return s.repo.GetEditorByEmail(ctx, email)
}

func (s *EditorService) Authenticate(ctx context.Context, email, password string) (any, error) {
	idToken, err := s.auth.VerifyPasswordWithREST(ctx, email, password)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	return map[string]string{
		"email":   email,
		"idToken": idToken,
	}, nil
}

func (s *EditorService) ChangePassword(ctx context.Context, email, newPassword string) error {
	user, err := s.auth.Client.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("firebase get user error: %w", err)
	}
	params := (&auth.UserToUpdate{}).Password(newPassword)
	_, err = s.auth.Client.UpdateUser(ctx, user.UID, params)
	if err != nil {
		return fmt.Errorf("firebase update user error: %w", err)
	}
	return nil
}
