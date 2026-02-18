package domain

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestAuthService_CreateUser(t *testing.T) {
	ctx := context.Background()
	email := "test@example.com"
	passwordHash := "hashed-pwd"
	userID := "user-123"
	token := "verify-token"

	repo := &MockRepository{
		CreateUserFunc: func(ctx context.Context, e, p string) (*User, error) {
			if e != email || p != passwordHash {
				return nil, errors.New("unexpected arguments to CreateUser")
			}
			return &User{ID: userID, Email: e}, nil
		},
		CreateEmailVerificationTokenFunc: func(ctx context.Context, t *EmailVerificationToken) error {
			if t.UserID != userID {
				return errors.New("unexpected userID in verification token")
			}
			t.Token = "hashed-token" // Simulate hashing in repo if needed, but service returns raw
			return nil
		},
	}

	publisher := &MockPublisher{
		PublishFunc: func(ctx context.Context, topic string, event interface{}) error {
			evt := event.(map[string]interface{})
			if evt["type"] != "user.registered" {
				t.Errorf("expected event type user.registered, got %v", evt["type"])
			}
			data := evt["data"].(map[string]string)
			if data["user_id"] != userID || data["email"] != email {
				t.Errorf("unexpected event data: %v", data)
			}
			if data["token"] == "" || !strings.Contains(data["link"], "/verify-email?token=") {
				t.Errorf("missing or invalid token/link in event data: %v", data)
			}
			return nil
		},
	}

	service := NewAuthService(repo, publisher)
	user, err := service.CreateUser(ctx, email, passwordHash)

	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	if user.ID != userID {
		t.Errorf("expected userID %s, got %s", userID, user.ID)
	}
}

func TestAuthService_CreatePasswordResetToken(t *testing.T) {
	ctx := context.Background()
	userID := "user-123"
	email := "test@example.com"

	repo := &MockRepository{
		GetUserByIDFunc: func(ctx context.Context, id string) (*User, error) {
			if id != userID {
				return nil, errors.New("user not found")
			}
			return &User{ID: userID, Email: email}, nil
		},
		CreatePasswordResetTokenFunc: func(ctx context.Context, token *PasswordResetToken) error {
			if token.UserID != userID {
				return errors.New("unexpected userID in reset token")
			}
			return nil
		},
	}

	publisher := &MockPublisher{
		PublishFunc: func(ctx context.Context, topic string, event interface{}) error {
			evt := event.(map[string]interface{})
			if evt["type"] != "password.reset" {
				t.Errorf("expected event type password.reset, got %v", evt["type"])
			}
			data := evt["data"].(map[string]string)
			if data["email"] != email || data["user_id"] != userID {
				t.Errorf("unexpected event data: %v", data)
			}
			return nil
		},
	}

	service := NewAuthService(repo, publisher)
	token, err := service.CreatePasswordResetToken(ctx, userID)

	if err != nil {
		t.Fatalf("CreatePasswordResetToken failed: %v", err)
	}
	if token == "" {
		t.Error("expected non-empty token")
	}
}

func TestAuthService_VerifyEmail(t *testing.T) {
	ctx := context.Background()
	rawToken := "raw-verify-token"
	userID := "user-123"

	repo := &MockRepository{
		GetEmailVerificationTokenFunc: func(ctx context.Context, hash string) (*EmailVerificationToken, error) {
			return &EmailVerificationToken{
				UserID: userID,
			}, nil
		},
		SetEmailVerifiedFunc: func(ctx context.Context, id string) error {
			if id != userID {
				return errors.New("unexpected userID in SetEmailVerified")
			}
			return nil
		},
		MarkEmailVerificationTokenUsedFunc: func(ctx context.Context, hash string) error {
			return nil
		},
	}

	service := NewAuthService(repo, nil)
	err := service.VerifyEmail(ctx, rawToken)

	if err != nil {
		t.Fatalf("VerifyEmail failed: %v", err)
	}
}
