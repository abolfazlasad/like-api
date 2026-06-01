package entities_test

import (
	"encoding/json"
	"testing"

	"like-api/internal/domain/entities"
)

func TestRole_Constants(t *testing.T) {
	if entities.RoleUser != "user" {
		t.Errorf("expected RoleUser='user', got %q", entities.RoleUser)
	}
	if entities.RoleAdmin != "admin" {
		t.Errorf("expected RoleAdmin='admin', got %q", entities.RoleAdmin)
	}
}

func TestUser_PasswordAccessible(t *testing.T) {
	u := entities.User{
		ID:       "u1",
		Username: "john",
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hashed-secret",
		Role:     entities.RoleUser,
	}

	if u.Password != "hashed-secret" {
		t.Error("password field should be accessible in Go code")
	}
}

func TestUser_PasswordNotSerialized(t *testing.T) {
	u := entities.User{
		ID:       "u1",
		Username: "alice",
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: "secret-hash",
		Role:     entities.RoleUser,
	}

	b, err := json.Marshal(u)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}

	if _, ok := m["password"]; ok {
		t.Error("password field must not appear in JSON output")
	}
}

func TestUser_JSONFields(t *testing.T) {
	u := entities.User{
		ID:       "u1",
		Username: "alice",
		Name:     "Alice",
		Email:    "alice@example.com",
		Role:     entities.RoleAdmin,
	}

	b, err := json.Marshal(u)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	expectedFields := []string{"id", "username", "name", "email", "role"}
	for _, field := range expectedFields {
		if _, ok := m[field]; !ok {
			t.Errorf("expected field %q in JSON", field)
		}
	}

	if m["role"] != "admin" {
		t.Errorf("expected role 'admin', got %v", m["role"])
	}
}

func TestUser_AdminRole(t *testing.T) {
	admin := entities.User{
		ID:   "a1",
		Role: entities.RoleAdmin,
	}

	if admin.Role != entities.RoleAdmin {
		t.Errorf("expected admin role, got %s", admin.Role)
	}
}
