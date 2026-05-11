package entities

// Like represents a like on a post by a user
type Like struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	PostID    string `json:"post_id"`
	CreatedAt string `json:"created_at"`
}
