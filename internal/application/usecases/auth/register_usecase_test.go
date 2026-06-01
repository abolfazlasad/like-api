package authusecase_test

import (
	"testing"

	"like-api/internal/application/dto/request"
	authusecase "like-api/internal/application/usecases/auth"
	"like-api/internal/infrastructure/database/memory"
)

const testSecret = "test-jwt-secret"

func TestRegister_Success(t *testing.T) {
	userRepo := memory.NewUserRepository()
	uc := authusecase.NewRegisterUseCase(userRepo, "test-secret")

	resp, err := uc.Execute(request.RegisterRequest{
		Username: "alice",
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: "password123",
	})

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp.Token == "" {
		t.Error("expected non-empty token")
	}
	if resp.User.Email != "alice@example.com" {
		t.Errorf("unexpected email: %s", resp.User.Email)
	}
	if resp.User.Role != "user" {
		t.Errorf("expected role=user, got %s", resp.User.Role)
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	userRepo := memory.NewUserRepository()
	uc := authusecase.NewRegisterUseCase(userRepo, "test-secret")

	req := request.RegisterRequest{
		Username: "alice",
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: "password123",
	}

	if _, err := uc.Execute(req); err != nil {
		t.Fatalf("first registration failed: %v", err)
	}

	req.Username = "alice2"
	_, err := uc.Execute(req)
	if err == nil {
		t.Fatal("expected error for duplicate email, got nil")
	}
	if err.Error() != "email already registered" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestRegister_TokenIsValid(t *testing.T) {
	userRepo := memory.NewUserRepository()
	uc := authusecase.NewRegisterUseCase(userRepo, "test-secret")

	resp, err := uc.Execute(request.RegisterRequest{
		Username: "bob",
		Name:     "Bob",
		Email:    "bob@example.com",
		Password: "password123",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Token) < 20 {
		t.Errorf("token looks too short: %s", resp.Token)
	}
}

func TestRegister_PasswordIsHashed(t *testing.T) {
	userRepo := memory.NewUserRepository()
	uc := authusecase.NewRegisterUseCase(userRepo, testSecret)

	req := request.RegisterRequest{
		Username: "bob",
		Name:     "Bob",
		Email:    "bob@example.com",
		Password: "plaintext",
	}
	resp, _ := uc.Execute(req)

	stored, _ := userRepo.FindByID(resp.User.ID)
	if stored.Password == "plaintext" {
		t.Error("password must be stored as bcrypt hash, not plaintext")
	}
	if len(stored.Password) < 20 {
		t.Error("stored password looks too short to be a bcrypt hash")
	}
}
