// Package interactions contains black-box tests for:
//   - POST /api/v1/videos/:id/like
//   - POST /api/v1/videos/:id/unlike
//   - POST /api/v1/videos/:id/view
package interactions

import (
	"fmt"
	"net/http"
	"testing"

	"like-api-tests/helpers"
)

// ─── Like ─────────────────────────────────────────────────────────────────────

func TestLike_Success(t *testing.T) {
	owner, _ := helpers.RegisterAndLogin(t, "like_owner")
	video := helpers.CreateVideo(t, owner, "Likeable Video")

	viewer, _ := helpers.RegisterAndLogin(t, "like_viewer")

	resp, err := helpers.DoRequest("POST", fmt.Sprintf("/api/v1/videos/%s/like", video.ID), nil, viewer)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusOK)

	var like helpers.Like
	apiResp := helpers.ParseResponse(t, resp, &like)

	if !apiResp.Success {
		t.Error("expected success=true")
	}
	if like.ID == "" {
		t.Error("expected non-empty like ID")
	}
	if like.PostID != video.ID {
		t.Errorf("post_id mismatch: want %s, got %s", video.ID, like.PostID)
	}
}

func TestLike_NoAuth(t *testing.T) {
	token, _ := helpers.RegisterAndLogin(t, "like_noauth_owner")
	video := helpers.CreateVideo(t, token, "Like No Auth Video")

	resp, err := helpers.DoRequest("POST", fmt.Sprintf("/api/v1/videos/%s/like", video.ID), nil, "")
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusUnauthorized)
	resp.Body.Close()
}

func TestLike_VideoNotFound(t *testing.T) {
	token, _ := helpers.RegisterAndLogin(t, "like_novid")

	resp, err := helpers.DoRequest("POST", "/api/v1/videos/00000000-0000-0000-0000-000000000000/like", nil, token)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusNotFound)
	resp.Body.Close()
}

func TestLike_Duplicate(t *testing.T) {
	owner, _ := helpers.RegisterAndLogin(t, "like_dup_owner")
	video := helpers.CreateVideo(t, owner, "Dup Like Video")
	viewer, _ := helpers.RegisterAndLogin(t, "like_dup_viewer")

	path := fmt.Sprintf("/api/v1/videos/%s/like", video.ID)

	// First like
	resp, _ := helpers.DoRequest("POST", path, nil, viewer)
	helpers.AssertStatus(t, resp, http.StatusOK)
	resp.Body.Close()

	// Second like — should conflict
	resp2, err := helpers.DoRequest("POST", path, nil, viewer)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp2, http.StatusConflict)
	resp2.Body.Close()
}

// ─── Unlike ───────────────────────────────────────────────────────────────────

func TestUnlike_Success(t *testing.T) {
	owner, _ := helpers.RegisterAndLogin(t, "unlike_owner")
	video := helpers.CreateVideo(t, owner, "Unlike Video")
	viewer, _ := helpers.RegisterAndLogin(t, "unlike_viewer")

	// Like first
	likeResp, _ := helpers.DoRequest("POST", fmt.Sprintf("/api/v1/videos/%s/like", video.ID), nil, viewer)
	helpers.AssertStatus(t, likeResp, http.StatusOK)
	likeResp.Body.Close()

	// Unlike
	resp, err := helpers.DoRequest("POST", fmt.Sprintf("/api/v1/videos/%s/unlike", video.ID), nil, viewer)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusOK)

	apiResp := helpers.ParseResponse(t, resp, nil)
	if !apiResp.Success {
		t.Error("expected success=true")
	}
}

func TestUnlike_NotPreviouslyLiked(t *testing.T) {
	owner, _ := helpers.RegisterAndLogin(t, "unlike_nolke_owner")
	video := helpers.CreateVideo(t, owner, "Unlike Never Liked")
	viewer, _ := helpers.RegisterAndLogin(t, "unlike_nolke_viewer")

	resp, err := helpers.DoRequest("POST", fmt.Sprintf("/api/v1/videos/%s/unlike", video.ID), nil, viewer)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusNotFound)
	resp.Body.Close()
}

func TestUnlike_NoAuth(t *testing.T) {
	token, _ := helpers.RegisterAndLogin(t, "unlike_noauth")
	video := helpers.CreateVideo(t, token, "Unlike No Auth")

	resp, err := helpers.DoRequest("POST", fmt.Sprintf("/api/v1/videos/%s/unlike", video.ID), nil, "")
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusUnauthorized)
	resp.Body.Close()
}

func TestUnlike_VideoNotFound(t *testing.T) {
	token, _ := helpers.RegisterAndLogin(t, "unlike_novid")

	resp, err := helpers.DoRequest("POST", "/api/v1/videos/00000000-0000-0000-0000-000000000000/unlike", nil, token)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusNotFound)
	resp.Body.Close()
}

// ─── View ─────────────────────────────────────────────────────────────────────

func TestView_AnonymousCounted(t *testing.T) {
	token, _ := helpers.RegisterAndLogin(t, "view_anon")
	video := helpers.CreateVideo(t, token, "Anon View Video")

	resp, err := helpers.DoRequest("POST", fmt.Sprintf("/api/v1/videos/%s/view", video.ID), nil, "")
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusOK)

	var viewData helpers.TrackViewResponse
	apiResp := helpers.ParseResponse(t, resp, &viewData)

	if !apiResp.Success {
		t.Error("expected success=true")
	}
	if !viewData.Counted {
		t.Error("expected counted=true for first view")
	}
}

func TestView_AuthenticatedCounted(t *testing.T) {
	owner, _ := helpers.RegisterAndLogin(t, "view_auth_owner")
	video := helpers.CreateVideo(t, owner, "Auth View Video")
	viewer, _ := helpers.RegisterAndLogin(t, "view_auth_viewer")

	resp, err := helpers.DoRequest("POST", fmt.Sprintf("/api/v1/videos/%s/view", video.ID), nil, viewer)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusOK)

	var viewData helpers.TrackViewResponse
	helpers.ParseResponse(t, resp, &viewData)

	if !viewData.Counted {
		t.Error("expected counted=true for first authenticated view")
	}
}

func TestView_DeduplicationWithinWindow(t *testing.T) {
	owner, _ := helpers.RegisterAndLogin(t, "view_dedup_owner")
	video := helpers.CreateVideo(t, owner, "Dedup View Video")
	viewer, _ := helpers.RegisterAndLogin(t, "view_dedup_viewer")

	path := fmt.Sprintf("/api/v1/videos/%s/view", video.ID)

	// First view — should be counted
	resp1, _ := helpers.DoRequest("POST", path, nil, viewer)
	helpers.AssertStatus(t, resp1, http.StatusOK)
	var v1 helpers.TrackViewResponse
	helpers.ParseResponse(t, resp1, &v1)
	if !v1.Counted {
		t.Error("first view: expected counted=true")
	}

	// Immediate second view — should NOT be counted (dedup window)
	resp2, err := helpers.DoRequest("POST", path, nil, viewer)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp2, http.StatusOK)
	var v2 helpers.TrackViewResponse
	helpers.ParseResponse(t, resp2, &v2)
	if v2.Counted {
		t.Error("second view within window: expected counted=false")
	}
}

func TestView_VideoNotFound(t *testing.T) {
	resp, err := helpers.DoRequest("POST", "/api/v1/videos/00000000-0000-0000-0000-000000000000/view", nil, "")
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusNotFound)
	resp.Body.Close()
}
