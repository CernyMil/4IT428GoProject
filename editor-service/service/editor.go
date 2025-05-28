package service

import (
	"context"
	"editor-service/models"
	uuidutil "editor-service/pkg/id"
	"editor-service/repository"
	transport "editor-service/transport/middleware"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	editor := &models.Editor{
		ID:             uuidutil.NewUUID(),
		Email:          userRecord.Email,
		FirstName:      firstName,
		LastName:       lastName,
		CreatedAt:      time.Now(),
		HashedPassword: string(hashedPassword),
	}
	return s.repo.CreateEditor(ctx, editor)
}

func (s *EditorService) GetByEmail(ctx context.Context, email string) (*models.Editor, error) {
	return s.repo.GetEditorByEmail(ctx, email)
}

/*
	func (s *EditorService) VerifyIDTokenAndGetEmail(ctx context.Context, idToken string) (string, error) {
		token, err := s.auth.Client.VerifyIDToken(ctx, idToken)
		if err != nil {
			return "", err
		}
		email, ok := token.Claims["email"].(string)
		if !ok || email == "" {
			return "", fmt.Errorf("email not found in token")
		}
		return email, nil
	}
*/
func (s *EditorService) Authenticate(ctx context.Context, email, password string) (*models.Editor, error) {
	editor, err := s.repo.GetEditorByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(editor.HashedPassword), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid password")
	}
	return editor, nil
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}
	if err := s.repo.UpdateEditorPassword(ctx, email, string(hashedPassword)); err != nil {
		return fmt.Errorf("failed to update local password: %w", err)
	}
	return nil
}
