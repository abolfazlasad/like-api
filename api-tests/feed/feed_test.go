// Package feed contains black-box tests for GET /api/v1/feed
package feed

import (
	"fmt"
	"net/http"
	"testing"

	"like-api-tests/helpers"
)

// ─── GET /api/v1/feed ─────────────────────────────────────────────────────────

func TestFeed_AnonymousDefaultPage(t *testing.T) {
	resp, err := helpers.DoRequest("GET", "/api/v1/feed", nil, "")
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusOK)

	var feed helpers.FeedResponse
	apiResp := helpers.ParseResponse(t, resp, &feed)

	if !apiResp.Success {
		t.Error("expected success=true")
	}
	if feed.Videos == nil {
		t.Error("expected videos array (can be empty but not nil)")
	}
}

func TestFeed_AuthenticatedRequest(t *testing.T) {
	token, _ := helpers.RegisterAndLogin(t, "feed_auth")

	resp, err := helpers.DoRequest("GET", "/api/v1/feed", nil, token)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusOK)
	resp.Body.Close()
}

func TestFeed_CustomLimit(t *testing.T) {
	resp, err := helpers.DoRequest("GET", "/api/v1/feed?limit=5", nil, "")
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusOK)

	var feed helpers.FeedResponse
	helpers.ParseResponse(t, resp, &feed)

	if len(feed.Videos) > 5 {
		t.Errorf("expected at most 5 videos with limit=5, got %d", len(feed.Videos))
	}
}

func TestFeed_Pagination(t *testing.T) {
	// Seed several videos so pagination can be exercised
	token, _ := helpers.RegisterAndLogin(t, "feed_page")
	for i := 0; i < 5; i++ {
		helpers.CreateVideo(t, token, fmt.Sprintf("Pagination Video %d", i))
	}

	// First page with limit=2
	resp, err := helpers.DoRequest("GET", "/api/v1/feed?limit=2", nil, "")
	if err != nil {
		t.Fatalf("page1 request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusOK)

	var page1 helpers.FeedResponse
	helpers.ParseResponse(t, resp, &page1)

	if len(page1.Videos) == 0 {
		t.Skip("no videos available to test pagination")
	}
	if page1.NextCursor == "" {
		// Only one page exists — still a valid scenario
		t.Log("only one page of results; skipping cursor follow-up")
		return
	}

	// Second page using cursor
	resp2, err := helpers.DoRequest("GET", "/api/v1/feed?limit=2&cursor="+page1.NextCursor, nil, "")
	if err != nil {
		t.Fatalf("page2 request error: %v", err)
	}
	helpers.AssertStatus(t, resp2, http.StatusOK)

	var page2 helpers.FeedResponse
	helpers.ParseResponse(t, resp2, &page2)

	// Verify no duplicate IDs across pages
	seen := map[string]bool{}
	for _, v := range page1.Videos {
		seen[v.ID] = true
	}
	for _, v := range page2.Videos {
		if seen[v.ID] {
			t.Errorf("duplicate video ID %s across pages", v.ID)
		}
	}
}

func TestFeed_InvalidCursor(t *testing.T) {
	resp, err := helpers.DoRequest("GET", "/api/v1/feed?cursor=not-a-timestamp", nil, "")
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusBadRequest)
	resp.Body.Close()
}

func TestFeed_LimitBoundaries(t *testing.T) {
	cases := []struct {
		limit          string
		expectOK       bool
		expectedStatus int
	}{
		{"1", true, http.StatusOK},
		{"100", true, http.StatusOK},
		{"0", false, http.StatusBadRequest},
		{"101", false, http.StatusBadRequest},
		{"-1", false, http.StatusBadRequest},
		{"abc", false, http.StatusBadRequest},
	}

	for _, tc := range cases {
		t.Run("limit="+tc.limit, func(t *testing.T) {
			resp, err := helpers.DoRequest("GET", "/api/v1/feed?limit="+tc.limit, nil, "")
			if err != nil {
				t.Fatalf("request error: %v", err)
			}
			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("limit=%s: expected %d, got %d", tc.limit, tc.expectedStatus, resp.StatusCode)
			}
			resp.Body.Close()
		})
	}
}

func TestFeed_NewestFirst(t *testing.T) {
	token, _ := helpers.RegisterAndLogin(t, "feed_order")
	helpers.CreateVideo(t, token, "Older Video")
	helpers.CreateVideo(t, token, "Newer Video")

	resp, err := helpers.DoRequest("GET", "/api/v1/feed?limit=10", nil, "")
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusOK)

	var feed helpers.FeedResponse
	helpers.ParseResponse(t, resp, &feed)

	// Verify descending order by created_at
	for i := 1; i < len(feed.Videos); i++ {
		if feed.Videos[i].CreatedAt.After(feed.Videos[i-1].CreatedAt) {
			t.Errorf("feed not sorted newest-first at index %d", i)
		}
	}
}
