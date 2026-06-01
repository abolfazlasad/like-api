// Package admin contains black-box tests for GET /api/v1/admin/users
//
// Note: these tests require an admin-role user to be available.
// The admin credentials are read from the environment variables:
//
//	ADMIN_EMAIL    (default: admin@example.com)
//	ADMIN_PASSWORD (default: AdminPass123!)
//
// If no admin account exists, the admin-specific tests are skipped gracefully
// (a 401/403 response confirms the endpoint enforces access control correctly).
package admin

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"like-api-tests/helpers"
)

// ─── helpers ──────────────────────────────────────────────────────────────────

func adminToken(t *testing.T) (string, bool) {
	t.Helper()
	email := os.Getenv("ADMIN_EMAIL")
	if email == "" {
		email = "admin@example.com"
	}
	password := os.Getenv("ADMIN_PASSWORD")
	if password == "" {
		password = "AdminPass123!"
	}

	resp, err := helpers.DoRequest("POST", "/api/v1/auth/login", helpers.LoginRequest{
		Email:    email,
		Password: password,
	}, "")
	if err != nil {
		t.Fatalf("admin login request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", false // no admin account available
	}

	var auth helpers.AuthResponse
	helpers.ParseResponse(t, resp, &auth)
	return auth.Token, true
}

// ─── GET /api/v1/admin/users ──────────────────────────────────────────────────

func TestAdminListUsers_WithAdminToken(t *testing.T) {
	token, ok := adminToken(t)
	if !ok {
		t.Skip("admin account not available — skipping admin-only test")
	}

	resp, err := helpers.DoRequest("GET", "/api/v1/admin/users", nil, token)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusOK)

	var users []helpers.UserDTO
	apiResp := helpers.ParseResponse(t, resp, &users)

	if !apiResp.Success {
		t.Error("expected success=true")
	}
	if users == nil {
		t.Error("expected non-nil users array")
	}

	// Every returned user must have required fields
	for i, u := range users {
		if u.ID == "" {
			t.Errorf("user[%d]: empty ID", i)
		}
		if u.Email == "" {
			t.Errorf("user[%d]: empty email", i)
		}
		if u.Role == "" {
			t.Errorf("user[%d]: empty role", i)
		}
	}
}

func TestAdminListUsers_NoToken(t *testing.T) {
	resp, err := helpers.DoRequest("GET", "/api/v1/admin/users", nil, "")
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusUnauthorized)
	resp.Body.Close()
}

func TestAdminListUsers_RegularUserToken(t *testing.T) {
	token, _ := helpers.RegisterAndLogin(t, "admin_regular")

	resp, err := helpers.DoRequest("GET", "/api/v1/admin/users", nil, token)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusForbidden)
	resp.Body.Close()
}

func TestAdminListUsers_InvalidToken(t *testing.T) {
	resp, err := helpers.DoRequest("GET", "/api/v1/admin/users", nil, "totally.invalid.token")
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusUnauthorized)
	resp.Body.Close()
}

func TestAdminListUsers_ResponseShape(t *testing.T) {
	token, ok := adminToken(t)
	if !ok {
		t.Skip("admin account not available")
	}

	resp, err := helpers.DoRequest("GET", "/api/v1/admin/users", nil, token)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusOK)

	// Verify top-level response envelope
	var raw map[string]json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		t.Fatalf("failed to decode raw response: %v", err)
	}
	resp.Body.Close()

	for _, field := range []string{"success", "data"} {
		if _, ok := raw[field]; !ok {
			t.Errorf("response missing field %q", field)
		}
	}
}
