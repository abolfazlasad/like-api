// Package videos contains black-box tests for:
//   - POST /api/v1/videos
//   - GET  /api/v1/videos/:id
package videos

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"like-api-tests/helpers"
)

// ─── POST /api/v1/videos ──────────────────────────────────────────────────────

func TestCreateVideo_Success(t *testing.T) {
	token, user := helpers.RegisterAndLogin(t, "vid_create")

	body := helpers.CreateVideoRequest{
		Title:       "My First Reel",
		VideoURL:    "https://cdn.example.com/video1.mp4",
		Description: "A test video",
	}
	resp, err := helpers.DoRequest("POST", "/api/v1/videos", body, token)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusCreated)

	var v helpers.Video
	apiResp := helpers.ParseResponse(t, resp, &v)

	if !apiResp.Success {
		t.Error("expected success=true")
	}
	if v.ID == "" {
		t.Error("expected non-empty video ID")
	}
	if v.Title != body.Title {
		t.Errorf("title mismatch: want %q, got %q", body.Title, v.Title)
	}
	if v.UserID != user.ID {
		t.Errorf("user_id mismatch: want %s, got %s", user.ID, v.UserID)
	}
}

func TestCreateVideo_NoAuth(t *testing.T) {
	body := helpers.CreateVideoRequest{
		Title:    "Unauthorized Video",
		VideoURL: "https://cdn.example.com/video.mp4",
	}
	resp, err := helpers.DoRequest("POST", "/api/v1/videos", body, "")
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusUnauthorized)
	resp.Body.Close()
}

func TestCreateVideo_InvalidToken(t *testing.T) {
	body := helpers.CreateVideoRequest{
		Title:    "Bad Token Video",
		VideoURL: "https://cdn.example.com/video.mp4",
	}
	resp, err := helpers.DoRequest("POST", "/api/v1/videos", body, "this.is.not.a.jwt")
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusUnauthorized)
	resp.Body.Close()
}

func TestCreateVideo_MissingTitle(t *testing.T) {
	token, _ := helpers.RegisterAndLogin(t, "vid_notitle")

	body := map[string]interface{}{
		"video_url": "https://cdn.example.com/video.mp4",
	}
	resp, err := helpers.DoRequest("POST", "/api/v1/videos", body, token)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusBadRequest)
	resp.Body.Close()
}

func TestCreateVideo_MissingVideoURL(t *testing.T) {
	token, _ := helpers.RegisterAndLogin(t, "vid_nourl")

	body := map[string]interface{}{
		"title": "No URL Video",
	}
	resp, err := helpers.DoRequest("POST", "/api/v1/videos", body, token)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusBadRequest)
	resp.Body.Close()
}

func TestCreateVideo_TitleTooLong(t *testing.T) {
	token, _ := helpers.RegisterAndLogin(t, "vid_longtitle")

	body := helpers.CreateVideoRequest{
		Title:    strings.Repeat("A", 256), // max is 255
		VideoURL: "https://cdn.example.com/video.mp4",
	}
	resp, err := helpers.DoRequest("POST", "/api/v1/videos", body, token)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusBadRequest)
	resp.Body.Close()
}

func TestCreateVideo_DescriptionTooLong(t *testing.T) {
	token, _ := helpers.RegisterAndLogin(t, "vid_longdesc")

	body := helpers.CreateVideoRequest{
		Title:       "Valid Title",
		VideoURL:    "https://cdn.example.com/video.mp4",
		Description: strings.Repeat("B", 2001), // max is 2000
	}
	resp, err := helpers.DoRequest("POST", "/api/v1/videos", body, token)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusBadRequest)
	resp.Body.Close()
}

// ─── GET /api/v1/videos/:id ───────────────────────────────────────────────────

func TestGetVideo_Success(t *testing.T) {
	token, _ := helpers.RegisterAndLogin(t, "vid_get")
	video := helpers.CreateVideo(t, token, "Get Test Video")

	resp, err := helpers.DoRequest("GET", fmt.Sprintf("/api/v1/videos/%s", video.ID), nil, "")
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusOK)

	var v helpers.Video
	helpers.ParseResponse(t, resp, &v)

	if v.ID != video.ID {
		t.Errorf("ID mismatch: want %s, got %s", video.ID, v.ID)
	}
	if v.Title != video.Title {
		t.Errorf("title mismatch: want %q, got %q", video.Title, v.Title)
	}
}

func TestGetVideo_WithAuth(t *testing.T) {
	token, _ := helpers.RegisterAndLogin(t, "vid_get_auth")
	video := helpers.CreateVideo(t, token, "Auth Get Test")

	// Authenticated request should also work
	resp, err := helpers.DoRequest("GET", fmt.Sprintf("/api/v1/videos/%s", video.ID), nil, token)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusOK)
	resp.Body.Close()
}

func TestGetVideo_NotFound(t *testing.T) {
	resp, err := helpers.DoRequest("GET", "/api/v1/videos/00000000-0000-0000-0000-000000000000", nil, "")
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusNotFound)
	resp.Body.Close()
}

func TestGetVideo_InvalidID(t *testing.T) {
	resp, err := helpers.DoRequest("GET", "/api/v1/videos/not-a-valid-uuid", nil, "")
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	// Either 400 or 404 is acceptable for a malformed UUID
	if resp.StatusCode != http.StatusBadRequest && resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 400 or 404, got %d", resp.StatusCode)
	}
	resp.Body.Close()
}
