package entities_test

import (
	"encoding/json"
	"testing"

	"like-api/internal/domain/entities"
)

func TestVideo_DefaultCounts(t *testing.T) {
	video := entities.Video{
		ID:       "v1",
		UserID:   "u1",
		Title:    "Test Video",
		VideoURL: "https://example.com/v.mp4",
	}

	if video.LikesCount != 0 {
		t.Errorf("expected LikesCount=0, got %d", video.LikesCount)
	}

	if video.ViewsCount != 0 {
		t.Errorf("expected ViewsCount=0, got %d", video.ViewsCount)
	}
}

func TestVideo_Fields(t *testing.T) {
	video := entities.Video{
		ID:          "v1",
		UserID:      "u1",
		Title:       "Go Architecture",
		Description: "A deep dive",
		VideoURL:    "https://example.com/v.mp4",
		LikesCount:  5,
		ViewsCount:  100,
		CreatedAt:   "2024-01-01T00:00:00Z",
	}

	if video.ID != "v1" {
		t.Errorf("expected ID=v1, got %s", video.ID)
	}

	if video.LikesCount != 5 {
		t.Errorf("expected LikesCount=5, got %d", video.LikesCount)
	}

	if video.ViewsCount != 100 {
		t.Errorf("expected ViewsCount=100, got %d", video.ViewsCount)
	}
}

func TestVideo_JSONFields(t *testing.T) {
	video := entities.Video{
		ID:          "v1",
		UserID:      "u1",
		Title:       "Test Video",
		Description: "desc",
		VideoURL:    "https://example.com/v.mp4",
		LikesCount:  5,
		ViewsCount:  100,
		CreatedAt:   "2024-01-01T00:00:00Z",
	}

	data, err := json.Marshal(video)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	expectedFields := []string{
		"id",
		"user_id",
		"title",
		"description",
		"video_url",
		"likes_count",
		"views_count",
		"created_at",
	}

	for _, field := range expectedFields {
		if _, exists := result[field]; !exists {
			t.Errorf("expected field %q in JSON", field)
		}
	}
}

func TestVideo_ZeroCounters(t *testing.T) {
	video := entities.Video{
		LikesCount: 0,
		ViewsCount: 0,
	}

	data, err := json.Marshal(video)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if result["likes_count"].(float64) != 0 {
		t.Error("expected likes_count to be 0")
	}

	if result["views_count"].(float64) != 0 {
		t.Error("expected views_count to be 0")
	}
}
