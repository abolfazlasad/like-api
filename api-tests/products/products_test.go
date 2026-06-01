// Package products contains black-box tests for:
//   - POST /api/v1/products
//   - GET  /api/v1/products/:id
//   - GET  /api/v1/videos/:id/product
package products

import (
	"fmt"
	"net/http"
	"testing"

	"like-api-tests/helpers"
)

// ─── POST /api/v1/products ────────────────────────────────────────────────────

func TestCreateProduct_Success(t *testing.T) {
	token, _ := helpers.RegisterAndLogin(t, "prod_create")
	video := helpers.CreateVideo(t, token, "Product Video")

	body := helpers.CreateProductRequest{
		Name:     "Cool Sneaker",
		Price:    99.99,
		ImageURL: "https://cdn.example.com/shoe.jpg",
		VideoID:  video.ID,
	}
	resp, err := helpers.DoRequest("POST", "/api/v1/products", body, token)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusCreated)

	var p helpers.Product
	apiResp := helpers.ParseResponse(t, resp, &p)

	if !apiResp.Success {
		t.Error("expected success=true")
	}
	if p.ID == "" {
		t.Error("expected non-empty product ID")
	}
	if p.VideoID != video.ID {
		t.Errorf("video_id mismatch: want %s, got %s", video.ID, p.VideoID)
	}
	if p.Price != body.Price {
		t.Errorf("price mismatch: want %f, got %f", body.Price, p.Price)
	}
}

func TestCreateProduct_NoAuth(t *testing.T) {
	body := helpers.CreateProductRequest{
		Name:     "Sneaker",
		Price:    50.0,
		ImageURL: "https://cdn.example.com/shoe.jpg",
		VideoID:  "00000000-0000-0000-0000-000000000001",
	}
	resp, err := helpers.DoRequest("POST", "/api/v1/products", body, "")
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusUnauthorized)
	resp.Body.Close()
}

func TestCreateProduct_VideoNotFound(t *testing.T) {
	token, _ := helpers.RegisterAndLogin(t, "prod_notvid")

	body := helpers.CreateProductRequest{
		Name:     "Ghost Product",
		Price:    10.0,
		ImageURL: "https://cdn.example.com/ghost.jpg",
		VideoID:  "00000000-0000-0000-0000-000000000000", // non-existent
	}
	resp, err := helpers.DoRequest("POST", "/api/v1/products", body, token)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusNotFound)
	resp.Body.Close()
}

func TestCreateProduct_DuplicateForSameVideo(t *testing.T) {
	token, _ := helpers.RegisterAndLogin(t, "prod_dup")
	video := helpers.CreateVideo(t, token, "Dup Product Video")

	body := helpers.CreateProductRequest{
		Name:     "First Product",
		Price:    20.0,
		ImageURL: "https://cdn.example.com/first.jpg",
		VideoID:  video.ID,
	}

	// First product — should succeed
	resp, _ := helpers.DoRequest("POST", "/api/v1/products", body, token)
	helpers.AssertStatus(t, resp, http.StatusCreated)
	resp.Body.Close()

	// Second product on same video — should conflict
	body.Name = "Second Product"
	resp2, err := helpers.DoRequest("POST", "/api/v1/products", body, token)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp2, http.StatusConflict)
	resp2.Body.Close()
}

func TestCreateProduct_MissingFields(t *testing.T) {
	token, _ := helpers.RegisterAndLogin(t, "prod_missing")
	video := helpers.CreateVideo(t, token, "Missing Fields Video")

	cases := []struct {
		name string
		body map[string]interface{}
	}{
		{
			"missing name",
			map[string]interface{}{"price": 10.0, "image_url": "https://cdn.example.com/img.jpg", "video_id": video.ID},
		},
		{
			"missing price",
			map[string]interface{}{"name": "X", "image_url": "https://cdn.example.com/img.jpg", "video_id": video.ID},
		},
		{
			"missing image_url",
			map[string]interface{}{"name": "X", "price": 10.0, "video_id": video.ID},
		},
		{
			"missing video_id",
			map[string]interface{}{"name": "X", "price": 10.0, "image_url": "https://cdn.example.com/img.jpg"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := helpers.DoRequest("POST", "/api/v1/products", tc.body, token)
			if err != nil {
				t.Fatalf("request error: %v", err)
			}
			helpers.AssertStatus(t, resp, http.StatusBadRequest)
			resp.Body.Close()
		})
	}
}

func TestCreateProduct_NameBoundaries(t *testing.T) {
	token, _ := helpers.RegisterAndLogin(t, "prod_namelen")
	video := helpers.CreateVideo(t, token, "Name Boundary Video")

	t.Run("empty name", func(t *testing.T) {
		body := helpers.CreateProductRequest{Name: "", Price: 1.0, ImageURL: "https://cdn.example.com/i.jpg", VideoID: video.ID}
		resp, _ := helpers.DoRequest("POST", "/api/v1/products", body, token)
		helpers.AssertStatus(t, resp, http.StatusBadRequest)
		resp.Body.Close()
	})

	t.Run("max-length name (255 chars)", func(t *testing.T) {
		video2 := helpers.CreateVideo(t, token, "Name 255 Video")
		body := helpers.CreateProductRequest{
			Name:     string(make([]byte, 255)),
			Price:    1.0,
			ImageURL: "https://cdn.example.com/i.jpg",
			VideoID:  video2.ID,
		}
		// Fill with 'A'
		for i := range body.Name {
			_ = i
		}
		body.Name = fmt.Sprintf("%255s", "A")
		resp, _ := helpers.DoRequest("POST", "/api/v1/products", body, token)
		// 201 or 400 depending on exact validation
		if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusBadRequest {
			t.Errorf("unexpected status %d", resp.StatusCode)
		}
		resp.Body.Close()
	})
}

// ─── GET /api/v1/products/:id ─────────────────────────────────────────────────

func TestGetProduct_Success(t *testing.T) {
	token, _ := helpers.RegisterAndLogin(t, "prod_get")
	video := helpers.CreateVideo(t, token, "Get Product Video")

	createBody := helpers.CreateProductRequest{
		Name:     "Gettable Product",
		Price:    49.99,
		ImageURL: "https://cdn.example.com/prod.jpg",
		VideoID:  video.ID,
	}
	createResp, _ := helpers.DoRequest("POST", "/api/v1/products", createBody, token)
	helpers.AssertStatus(t, createResp, http.StatusCreated)
	var p helpers.Product
	helpers.ParseResponse(t, createResp, &p)

	// Now fetch by ID
	resp, err := helpers.DoRequest("GET", fmt.Sprintf("/api/v1/products/%s", p.ID), nil, "")
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusOK)

	var fetched helpers.Product
	helpers.ParseResponse(t, resp, &fetched)

	if fetched.ID != p.ID {
		t.Errorf("ID mismatch: want %s, got %s", p.ID, fetched.ID)
	}
}

func TestGetProduct_NotFound(t *testing.T) {
	resp, err := helpers.DoRequest("GET", "/api/v1/products/00000000-0000-0000-0000-000000000000", nil, "")
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusNotFound)
	resp.Body.Close()
}

// ─── GET /api/v1/videos/:id/product ──────────────────────────────────────────

func TestGetVideoProduct_Success(t *testing.T) {
	token, _ := helpers.RegisterAndLogin(t, "vidprod_get")
	video := helpers.CreateVideo(t, token, "Video With Product")

	createBody := helpers.CreateProductRequest{
		Name:     "Linked Product",
		Price:    75.0,
		ImageURL: "https://cdn.example.com/linked.jpg",
		VideoID:  video.ID,
	}
	createResp, _ := helpers.DoRequest("POST", "/api/v1/products", createBody, token)
	helpers.AssertStatus(t, createResp, http.StatusCreated)
	createResp.Body.Close()

	resp, err := helpers.DoRequest("GET", fmt.Sprintf("/api/v1/videos/%s/product", video.ID), nil, "")
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusOK)

	var p helpers.Product
	helpers.ParseResponse(t, resp, &p)

	if p.VideoID != video.ID {
		t.Errorf("video_id mismatch: want %s, got %s", video.ID, p.VideoID)
	}
}

func TestGetVideoProduct_VideoHasNoProduct(t *testing.T) {
	token, _ := helpers.RegisterAndLogin(t, "vidprod_none")
	video := helpers.CreateVideo(t, token, "Video Without Product")

	resp, err := helpers.DoRequest("GET", fmt.Sprintf("/api/v1/videos/%s/product", video.ID), nil, "")
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusNotFound)
	resp.Body.Close()
}

func TestGetVideoProduct_VideoNotFound(t *testing.T) {
	resp, err := helpers.DoRequest("GET", "/api/v1/videos/00000000-0000-0000-0000-000000000000/product", nil, "")
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	helpers.AssertStatus(t, resp, http.StatusNotFound)
	resp.Body.Close()
}
