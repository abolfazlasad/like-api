package entities

// Post represents a post in the system
type Post struct {
	ID      string `json:"id"`
	UserID  string `json:"user_id"`
	Content string `json:"content"`
	Title   string `json:"title"`
	Likes   int    `json:"likes"`
}
