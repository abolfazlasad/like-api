package entities

// User represents a user in the system
type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"-"` // bcrypt hash — never serialized
	CreatedAt string `json:"created_at"`
}
