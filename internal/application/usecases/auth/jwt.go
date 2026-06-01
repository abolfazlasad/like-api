package authusecase

import (
	"time"

	"like-api/internal/domain/entities"

	"github.com/golang-jwt/jwt/v5"
)

// Claims holds the JWT payload including the user's role.
type Claims struct {
	UserID   string        `json:"user_id"`
	Username string        `json:"username"`
	Role     entities.Role `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(user entities.User, secret string) (string, error) {
	role := user.Role
	if role == "" {
		role = entities.RoleUser
	}

	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}
