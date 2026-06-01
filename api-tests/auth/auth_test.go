// Package auth contains black-box tests for POST /api/v1/auth/register
// and POST /api/v1/auth/login.
package auth

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"like-api-tests/helpers"
)

// ─── Register ────────────────────────────────────────────────────────────────

func TestRegister_Success(t *testing.T) {
	ts := time.Now().UnixNano()
	body := helpers.RegisterRequest{
		Email:    fmt.Sprintf("newuser_%d@example.com", ts),
		Name:     "New User",
		Password: "SecurePass1!",
		Username: fmt.Sprintf("newuser%d", ts%1_000_000),
	}

	resp, err := helpers.DoRequest("POST", "/api/v1/auth/register", body, "")
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusCreated)

	var auth helpers.AuthResponse
	apiResp := helpers.ParseResponse(t, resp, &auth)

	if !apiResp.Success {
		t.Errorf("expected success=true, got false")
	}
	if auth.Token == "" {
		t.Error("expected non-empty token")
	}
	if auth.User.Email != body.Email {
		t.Errorf("email mismatch: want %s, got %s", body.Email, auth.User.Email)
	}
	if auth.User.Role != "user" {
		t.Errorf("expected role=user, got %s", auth.User.Role)
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	ts := time.Now().UnixNano()
	body := helpers.RegisterRequest{
		Email:    fmt.Sprintf("dup_%d@example.com", ts),
		Name:     "Dup User",
		Password: "SecurePass1!",
		Username: fmt.Sprintf("dup%d", ts%1_000_000),
	}

	// First registration — should succeed
	resp, _ := helpers.DoRequest("POST", "/api/v1/auth/register", body, "")
	helpers.AssertStatus(t, resp, http.StatusCreated)
	resp.Body.Close()

	// Second registration with same email — should conflict
	body.Username = fmt.Sprintf("dup2_%d", ts%1_000_000)
	resp2, err := helpers.DoRequest("POST", "/api/v1/auth/register", body, "")
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp2, http.StatusConflict)

	apiResp := helpers.ParseResponse(t, resp2, nil)
	if apiResp.Success {
		t.Error("expected success=false for duplicate email")
	}
}

func TestRegister_MissingRequiredFields(t *testing.T) {
	cases := []struct {
		name string
		body map[string]interface{}
	}{
		{
			name: "missing email",
			body: map[string]interface{}{"name": "A", "password": "Pass1234!", "username": "abc123"},
		},
		{
			name: "missing password",
			body: map[string]interface{}{"email": "a@b.com", "name": "A", "username": "abc123"},
		},
		{
			name: "missing username",
			body: map[string]interface{}{"email": "a@b.com", "name": "A", "password": "Pass1234!"},
		},
		{
			name: "missing name",
			body: map[string]interface{}{"email": "a@b.com", "password": "Pass1234!", "username": "abc123"},
		},
		{
			name: "empty body",
			body: map[string]interface{}{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := helpers.DoRequest("POST", "/api/v1/auth/register", tc.body, "")
			if err != nil {
				t.Fatalf("request error: %v", err)
			}
			helpers.AssertStatus(t, resp, http.StatusBadRequest)
			resp.Body.Close()
		})
	}
}

func TestRegister_PasswordTooShort(t *testing.T) {
	ts := time.Now().UnixNano()
	body := helpers.RegisterRequest{
		Email:    fmt.Sprintf("short_%d@example.com", ts),
		Name:     "Short Pass",
		Password: "1234567", // 7 chars, minimum is 8
		Username: fmt.Sprintf("short%d", ts%1_000_000),
	}
	resp, err := helpers.DoRequest("POST", "/api/v1/auth/register", body, "")
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusBadRequest)
	resp.Body.Close()
}

func TestRegister_UsernameTooShort(t *testing.T) {
	ts := time.Now().UnixNano()
	body := helpers.RegisterRequest{
		Email:    fmt.Sprintf("un_%d@example.com", ts),
		Name:     "Short UN",
		Password: "SecurePass1!",
		Username: "ab", // min is 3
	}
	resp, err := helpers.DoRequest("POST", "/api/v1/auth/register", body, "")
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusBadRequest)
	resp.Body.Close()
}

func TestRegister_UsernameTooLong(t *testing.T) {
	ts := time.Now().UnixNano()
	body := helpers.RegisterRequest{
		Email:    fmt.Sprintf("ulong_%d@example.com", ts),
		Name:     "Long UN",
		Password: "SecurePass1!",
		Username: "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxy", // 51 chars, max is 50
	}
	resp, err := helpers.DoRequest("POST", "/api/v1/auth/register", body, "")
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusBadRequest)
	resp.Body.Close()
}

// ─── Login ───────────────────────────────────────────────────────────────────

func TestLogin_Success(t *testing.T) {
	ts := time.Now().UnixNano()
	email := fmt.Sprintf("login_%d@example.com", ts)
	password := "LoginPass1!"

	// Register first
	regBody := helpers.RegisterRequest{
		Email:    email,
		Name:     "Login User",
		Password: password,
		Username: fmt.Sprintf("loginuser%d", ts%1_000_000),
	}
	resp, _ := helpers.DoRequest("POST", "/api/v1/auth/register", regBody, "")
	helpers.AssertStatus(t, resp, http.StatusCreated)
	resp.Body.Close()

	// Login
	loginBody := helpers.LoginRequest{Email: email, Password: password}
	resp2, err := helpers.DoRequest("POST", "/api/v1/auth/login", loginBody, "")
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp2, http.StatusOK)

	var auth helpers.AuthResponse
	apiResp := helpers.ParseResponse(t, resp2, &auth)

	if !apiResp.Success {
		t.Error("expected success=true")
	}
	if auth.Token == "" {
		t.Error("expected non-empty token")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	ts := time.Now().UnixNano()
	email := fmt.Sprintf("wrongpw_%d@example.com", ts)

	regBody := helpers.RegisterRequest{
		Email:    email,
		Name:     "Wrong PW",
		Password: "CorrectPass1!",
		Username: fmt.Sprintf("wrongpw%d", ts%1_000_000),
	}
	resp, _ := helpers.DoRequest("POST", "/api/v1/auth/register", regBody, "")
	helpers.AssertStatus(t, resp, http.StatusCreated)
	resp.Body.Close()

	loginBody := helpers.LoginRequest{Email: email, Password: "WrongPassword!"}
	resp2, err := helpers.DoRequest("POST", "/api/v1/auth/login", loginBody, "")
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp2, http.StatusUnauthorized)
	resp2.Body.Close()
}

func TestLogin_UnregisteredEmail(t *testing.T) {
	loginBody := helpers.LoginRequest{
		Email:    "nobody_exists@example.com",
		Password: "SomePass1!",
	}
	resp, err := helpers.DoRequest("POST", "/api/v1/auth/login", loginBody, "")
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusUnauthorized)
	resp.Body.Close()
}

func TestLogin_MissingFields(t *testing.T) {
	cases := []struct {
		name string
		body map[string]interface{}
	}{
		{"missing email", map[string]interface{}{"password": "Pass1234!"}},
		{"missing password", map[string]interface{}{"email": "a@b.com"}},
		{"empty body", map[string]interface{}{}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := helpers.DoRequest("POST", "/api/v1/auth/login", tc.body, "")
			if err != nil {
				t.Fatalf("request error: %v", err)
			}
			helpers.AssertStatus(t, resp, http.StatusBadRequest)
			resp.Body.Close()
		})
	}
}
