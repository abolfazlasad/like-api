package entities_test

import (
	"encoding/json"
	"testing"

	"like-api/internal/domain/entities"
)

func TestLike_JSONFields(t *testing.T) {
	l := entities.Like{
		ID:        "l1",
		UserID:    "u1",
		PostID:    "v1",
		CreatedAt: "2024-01-01T00:00:00Z",
	}

	b, err := json.Marshal(l)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var m map[string]interface{}
	_ = json.Unmarshal(b, &m)

	for _, field := range []string{"id", "user_id", "post_id", "created_at"} {
		if _, ok := m[field]; !ok {
			t.Errorf("expected field %q in JSON", field)
		}
	}
}
