package authusecase_test

import (
	"testing"
	"time"

	authusecase "like-api/internal/application/usecases/auth"
	"like-api/internal/domain/entities"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateToken_ReturnsValidToken(t *testing.T) {
	secret := "test-secret"

	user := entities.User{
		ID:       "user-1",
		Username: "john",
		Role:     entities.RoleAdmin,
	}

	tokenString, err := authusecase.GenerateToken(user, secret)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&authusecase.Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	)
	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}

	claims, ok := token.Claims.(*authusecase.Claims)
	if !ok {
		t.Fatal("expected Claims type")
	}

	if claims.UserID != user.ID {
		t.Errorf("expected UserID=%s, got %s", user.ID, claims.UserID)
	}

	if claims.Username != user.Username {
		t.Errorf(
			"expected Username=%s, got %s",
			user.Username,
			claims.Username,
		)
	}

	if claims.Role != entities.RoleAdmin {
		t.Errorf(
			"expected Role=%s, got %s",
			entities.RoleAdmin,
			claims.Role,
		)
	}

	if claims.ExpiresAt == nil {
		t.Fatal("expected ExpiresAt to be set")
	}

	if claims.IssuedAt == nil {
		t.Fatal("expected IssuedAt to be set")
	}

	lifetime := claims.ExpiresAt.Time.Sub(claims.IssuedAt.Time)

	if lifetime < 23*time.Hour || lifetime > 25*time.Hour {
		t.Errorf(
			"expected lifetime around 24h, got %v",
			lifetime,
		)
	}
}

func TestGenerateToken_DefaultsToUserRole(t *testing.T) {
	secret := "test-secret"

	user := entities.User{
		ID:       "user-1",
		Username: "john",
		Role:     "",
	}

	tokenString, err := authusecase.GenerateToken(user, secret)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&authusecase.Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	)
	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}

	claims := token.Claims.(*authusecase.Claims)

	if claims.Role != entities.RoleUser {
		t.Errorf(
			"expected Role=%s, got %s",
			entities.RoleUser,
			claims.Role,
		)
	}
}
