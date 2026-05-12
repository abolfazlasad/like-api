package request

// LikePostRequest represents the request body for liking a post
type LikePostRequest struct {
	UserID string `json:"user_id" binding:"required" example:"user1"`
	PostID string `json:"post_id" binding:"required" example:"post1"`
}

// UnlikePostRequest represents the request body for unliking a post
type UnlikePostRequest struct {
	UserID string `json:"user_id" binding:"required" example:"user1"`
	PostID string `json:"post_id" binding:"required" example:"post1"`
}
