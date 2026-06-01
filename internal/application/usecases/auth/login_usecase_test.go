package authusecase_test

import (
	"testing"

	"like-api/internal/application/dto/request"
	authusecase "like-api/internal/application/usecases/auth"
	"like-api/internal/domain/repositories"
	"like-api/internal/infrastructure/database/memory"
)

func registerUser(t *testing.T, userRepo repositories.UserRepository, email, password string) {
	t.Helper()
	uc := authusecase.NewRegisterUseCase(userRepo, "test-secret")
	_, err := uc.Execute(request.RegisterRequest{
		Username: "testuser",
		Name:     "Test User",
		Email:    email,
		Password: password,
	})
	if err != nil {
		t.Fatalf("setup register failed: %v", err)
	}
}

func TestLogin_Success(t *testing.T) {
	userRepo := memory.NewUserRepository()
	registerUser(t, userRepo, "user@example.com", "mypassword")

	loginUC := authusecase.NewLoginUseCase(userRepo, "test-secret")
	resp, err := loginUC.Execute(request.LoginRequest{
		Email:    "user@example.com",
		Password: "mypassword",
	})

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp.Token == "" {
		t.Error("expected non-empty token")
	}
	if resp.User.Email != "user@example.com" {
		t.Errorf("unexpected email in response: %s", resp.User.Email)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	userRepo := memory.NewUserRepository()
	registerUser(t, userRepo, "user@example.com", "correctpassword")

	loginUC := authusecase.NewLoginUseCase(userRepo, "test-secret")
	_, err := loginUC.Execute(request.LoginRequest{
		Email:    "user@example.com",
		Password: "wrongpassword",
	})

	if err == nil {
		t.Fatal("expected error for wrong password, got nil")
	}
	if err.Error() != "invalid credentials" {
		t.Errorf("unexpected error: %s", err.Error())
	}
}

func TestLogin_UnknownEmail(t *testing.T) {
	userRepo := memory.NewUserRepository()

	loginUC := authusecase.NewLoginUseCase(userRepo, "test-secret")
	_, err := loginUC.Execute(request.LoginRequest{
		Email:    "nobody@example.com",
		Password: "any",
	})

	if err == nil {
		t.Fatal("expected error for unknown email, got nil")
	}
	if err.Error() != "invalid credentials" {
		t.Errorf("unexpected error: %s", err.Error())
	}
}

func TestLogin_ErrorMessageObfuscated(t *testing.T) {
	userRepo := memory.NewUserRepository()

	registerUC := authusecase.NewRegisterUseCase(userRepo, testSecret)
	loginUC := authusecase.NewLoginUseCase(userRepo, testSecret)

	_, err := registerUC.Execute(request.RegisterRequest{
		Username: "eve",
		Name:     "Eve",
		Email:    "eve@example.com",
		Password: "correctpass",
	})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	_, errWrongPass := loginUC.Execute(request.LoginRequest{
		Email:    "eve@example.com",
		Password: "wrong",
	})

	_, errNoUser := loginUC.Execute(request.LoginRequest{
		Email:    "ghost@example.com",
		Password: "wrong",
	})

	if errWrongPass == nil {
		t.Fatal("expected error for wrong password")
	}

	if errNoUser == nil {
		t.Fatal("expected error for non-existent user")
	}

	if errWrongPass.Error() != errNoUser.Error() {
		t.Fatalf(
			"expected same error message, got %q and %q",
			errWrongPass.Error(),
			errNoUser.Error(),
		)
	}
}
