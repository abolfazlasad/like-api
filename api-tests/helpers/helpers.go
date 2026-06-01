package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

var BaseURL = os.Getenv("BASE_URL")

// ─── Generic response wrapper ───────────────────────────────────────────────

type APIResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

// ─── Auth DTOs ───────────────────────────────────────────────────────────────

type RegisterRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string  `json:"token"`
	User  UserDTO `json:"user"`
}

type UserDTO struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

// ─── Video DTOs ──────────────────────────────────────────────────────────────

type CreateVideoRequest struct {
	Title       string `json:"title"`
	VideoURL    string `json:"video_url"`
	Description string `json:"description,omitempty"`
}

type Video struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	VideoURL    string    `json:"video_url"`
	UserID      string    `json:"user_id"`
	LikesCount  int       `json:"likes_count"`
	ViewsCount  int       `json:"views_count"`
	CreatedAt   time.Time `json:"created_at"`
}

// ─── Product DTOs ─────────────────────────────────────────────────────────────

type CreateProductRequest struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	ImageURL string  `json:"image_url"`
	VideoID  string  `json:"video_id"`
}

type Product struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	ImageURL string  `json:"image_url"`
	VideoID  string  `json:"video_id"`
}

// ─── Other DTOs ───────────────────────────────────────────────────────────────

type Like struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	PostID    string    `json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
}

type TrackViewResponse struct {
	Counted bool `json:"counted"`
}

type VideoStats struct {
	VideoID        string  `json:"video_id"`
	Views          int     `json:"views"`
	Likes          int     `json:"likes"`
	EngagementRate float64 `json:"engagement_rate"`
}

type FeedResponse struct {
	Videos     []Video `json:"videos"`
	NextCursor string  `json:"nextCursor"`
}

// ─── HTTP helpers ─────────────────────────────────────────────────────────────

func DoRequest(method, path string, body interface{}, token string) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, BaseURL+path, bodyReader)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	return client.Do(req)
}

func ParseResponse(t *testing.T, resp *http.Response, out interface{}) APIResponse {
	t.Helper()
	defer resp.Body.Close()
	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if out != nil && apiResp.Data != nil {
		if err := json.Unmarshal(apiResp.Data, out); err != nil {
			t.Fatalf("failed to decode data field: %v", err)
		}
	}
	return apiResp
}

func AssertStatus(t *testing.T, resp *http.Response, expected int) {
	t.Helper()
	if resp.StatusCode != expected {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected status %d, got %d — body: %s", expected, resp.StatusCode, string(body))
	}
}

// ─── Seed helpers ─────────────────────────────────────────────────────────────

// RegisterAndLogin registers a unique user and returns token + user info.
func RegisterAndLogin(t *testing.T, suffix string) (token string, user UserDTO) {
	t.Helper()
	ts := time.Now().UnixNano()
	reg := RegisterRequest{
		Email:    fmt.Sprintf("user_%s_%d@test.com", suffix, ts),
		Name:     fmt.Sprintf("Test %s", suffix),
		Password: "Password123!",
		Username: fmt.Sprintf("user_%s_%d", suffix, ts%1_000_000),
	}
	resp, err := DoRequest("POST", "/api/v1/auth/register", reg, "")
	if err != nil {
		t.Fatalf("register request failed: %v", err)
	}
	AssertStatus(t, resp, 201)

	var auth AuthResponse
	ParseResponse(t, resp, &auth)
	return auth.Token, auth.User
}

// CreateVideo creates a video and returns it.
func CreateVideo(t *testing.T, token string, title string) Video {
	t.Helper()
	req := CreateVideoRequest{
		Title:    title,
		VideoURL: "https://example.com/video.mp4",
	}
	resp, err := DoRequest("POST", "/api/v1/videos", req, token)
	if err != nil {
		t.Fatalf("create video request failed: %v", err)
	}
	AssertStatus(t, resp, 201)

	var v Video
	ParseResponse(t, resp, &v)
	return v
}
