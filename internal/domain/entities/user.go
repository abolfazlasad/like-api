package entities

// Role defines the access level of a user.
type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

// User represents a user in the system.
type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"-"` // bcrypt hash — never serialized
	Role      Role   `json:"role"`
	CreatedAt string `json:"created_at"`
}
