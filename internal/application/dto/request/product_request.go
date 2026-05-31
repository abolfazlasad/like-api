package request

// CreateProductRequest represents the request body for creating a product
type CreateProductRequest struct {
	VideoID  string  `json:"video_id"  binding:"required"`
	Name     string  `json:"name"      binding:"required,min=1,max=255"`
	Price    float64 `json:"price"     binding:"required,gt=0"`
	ImageURL string  `json:"image_url" binding:"required,url"`
}
