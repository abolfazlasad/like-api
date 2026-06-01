// Package analytics contains black-box tests for GET /api/v1/videos/:id/stats
package analytics

import (
	"fmt"
	"math"
	"net/http"
	"testing"

	"like-api-tests/helpers"
)

// ─── GET /api/v1/videos/:id/stats ────────────────────────────────────────────

func TestStats_ZeroCountsOnNewVideo(t *testing.T) {
	token, _ := helpers.RegisterAndLogin(t, "stats_zero")
	video := helpers.CreateVideo(t, token, "Zero Stats Video")

	resp, err := helpers.DoRequest("GET", fmt.Sprintf("/api/v1/videos/%s/stats", video.ID), nil, "")
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusOK)

	var stats helpers.VideoStats
	apiResp := helpers.ParseResponse(t, resp, &stats)

	if !apiResp.Success {
		t.Error("expected success=true")
	}
	if stats.VideoID != video.ID {
		t.Errorf("video_id mismatch: want %s, got %s", video.ID, stats.VideoID)
	}
	if stats.Views != 0 {
		t.Errorf("expected 0 views, got %d", stats.Views)
	}
	if stats.Likes != 0 {
		t.Errorf("expected 0 likes, got %d", stats.Likes)
	}
	if stats.EngagementRate != 0 {
		t.Errorf("expected 0 engagement_rate when views=0, got %f", stats.EngagementRate)
	}
}

func TestStats_AfterView(t *testing.T) {
	owner, _ := helpers.RegisterAndLogin(t, "stats_view_owner")
	video := helpers.CreateVideo(t, owner, "View Stats Video")
	viewer, _ := helpers.RegisterAndLogin(t, "stats_view_viewer")

	// Record a view
	viewResp, _ := helpers.DoRequest("POST", fmt.Sprintf("/api/v1/videos/%s/view", video.ID), nil, viewer)
	helpers.AssertStatus(t, viewResp, http.StatusOK)
	viewResp.Body.Close()

	resp, err := helpers.DoRequest("GET", fmt.Sprintf("/api/v1/videos/%s/stats", video.ID), nil, "")
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusOK)

	var stats helpers.VideoStats
	helpers.ParseResponse(t, resp, &stats)

	if stats.Views != 1 {
		t.Errorf("expected 1 view, got %d", stats.Views)
	}
	if stats.EngagementRate != 0 {
		t.Errorf("expected 0 engagement_rate (no likes), got %f", stats.EngagementRate)
	}
}

func TestStats_AfterLikeAndView(t *testing.T) {
	owner, _ := helpers.RegisterAndLogin(t, "stats_both_owner")
	video := helpers.CreateVideo(t, owner, "Full Stats Video")
	viewer, _ := helpers.RegisterAndLogin(t, "stats_both_viewer")

	// Record a view
	viewResp, _ := helpers.DoRequest("POST", fmt.Sprintf("/api/v1/videos/%s/view", video.ID), nil, viewer)
	helpers.AssertStatus(t, viewResp, http.StatusOK)
	viewResp.Body.Close()

	// Like the video
	likeResp, _ := helpers.DoRequest("POST", fmt.Sprintf("/api/v1/videos/%s/like", video.ID), nil, viewer)
	helpers.AssertStatus(t, likeResp, http.StatusOK)
	likeResp.Body.Close()

	resp, err := helpers.DoRequest("GET", fmt.Sprintf("/api/v1/videos/%s/stats", video.ID), nil, "")
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusOK)

	var stats helpers.VideoStats
	helpers.ParseResponse(t, resp, &stats)

	if stats.Views < 1 {
		t.Errorf("expected at least 1 view, got %d", stats.Views)
	}
	if stats.Likes < 1 {
		t.Errorf("expected at least 1 like, got %d", stats.Likes)
	}

	// engagement_rate = (likes / views) * 100
	expectedRate := (float64(stats.Likes) / float64(stats.Views)) * 100
	if math.Abs(stats.EngagementRate-expectedRate) > 0.01 {
		t.Errorf("engagement_rate: want %.4f, got %.4f", expectedRate, stats.EngagementRate)
	}
}

func TestStats_EngagementRateZeroWhenNoViews(t *testing.T) {
	// Create a video and like it without viewing — edge case
	owner, _ := helpers.RegisterAndLogin(t, "stats_eng_owner")
	video := helpers.CreateVideo(t, owner, "Liked But Not Viewed")
	liker, _ := helpers.RegisterAndLogin(t, "stats_eng_liker")

	likeResp, _ := helpers.DoRequest("POST", fmt.Sprintf("/api/v1/videos/%s/like", video.ID), nil, liker)
	helpers.AssertStatus(t, likeResp, http.StatusOK)
	likeResp.Body.Close()

	resp, err := helpers.DoRequest("GET", fmt.Sprintf("/api/v1/videos/%s/stats", video.ID), nil, "")
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusOK)

	var stats helpers.VideoStats
	helpers.ParseResponse(t, resp, &stats)

	// API spec: returns 0 when views = 0
	if stats.Views == 0 && stats.EngagementRate != 0 {
		t.Errorf("expected engagement_rate=0 when views=0, got %f", stats.EngagementRate)
	}
}

func TestStats_AccessibleAnonymously(t *testing.T) {
	token, _ := helpers.RegisterAndLogin(t, "stats_anon")
	video := helpers.CreateVideo(t, token, "Anon Stats Video")

	resp, err := helpers.DoRequest("GET", fmt.Sprintf("/api/v1/videos/%s/stats", video.ID), nil, "")
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusOK)
	resp.Body.Close()
}

func TestStats_VideoNotFound(t *testing.T) {
	resp, err := helpers.DoRequest("GET", "/api/v1/videos/00000000-0000-0000-0000-000000000000/stats", nil, "")
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusNotFound)
	resp.Body.Close()
}

func TestStats_WithOptionalAuth(t *testing.T) {
	owner, _ := helpers.RegisterAndLogin(t, "stats_optauth_own")
	video := helpers.CreateVideo(t, owner, "Optional Auth Stats")
	viewer, _ := helpers.RegisterAndLogin(t, "stats_optauth_view")

	// Authenticated request should also work
	resp, err := helpers.DoRequest("GET", fmt.Sprintf("/api/v1/videos/%s/stats", video.ID), nil, viewer)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusOK)
	resp.Body.Close()
}
