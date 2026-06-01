package entities_test

import (
	"encoding/json"
	"testing"

	"like-api/internal/domain/entities"
)

func TestProduct_JSONFields(t *testing.T) {
	p := entities.Product{
		ID:       "p1",
		VideoID:  "v1",
		Name:     "Test Product",
		Price:    29.99,
		ImageURL: "https://example.com/img.jpg",
	}

	b, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var m map[string]interface{}
	_ = json.Unmarshal(b, &m)

	for _, field := range []string{"id", "video_id", "name", "price", "image_url"} {
		if _, ok := m[field]; !ok {
			t.Errorf("expected field %q in JSON", field)
		}
	}

	if m["price"].(float64) != 29.99 {
		t.Errorf("expected price 29.99, got %v", m["price"])
	}
}
