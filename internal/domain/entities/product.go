package entities

// Product represents a shoppable product attached to a video
type Product struct {
	ID       string  `json:"id"`
	VideoID  string  `json:"video_id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	ImageURL string  `json:"image_url"`
}
