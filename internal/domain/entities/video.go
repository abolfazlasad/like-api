package entities

// Video represents a video/reel uploaded by a user
type Video struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	VideoURL    string `json:"video_url"`
	LikesCount  int    `json:"likes_count"`
	ViewsCount  int    `json:"views_count"`
	CreatedAt   string `json:"created_at"`
}
