package request

// CreateVideoRequest represents the request body for uploading a video
type CreateVideoRequest struct {
	Title       string `json:"title"       binding:"required,min=1,max=255"`
	Description string `json:"description" binding:"max=2000"`
	VideoURL    string `json:"video_url"   binding:"required,url"`
}
