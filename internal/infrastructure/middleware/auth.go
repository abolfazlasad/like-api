package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	authusecase "like-api/internal/application/usecases/auth"
	"like-api/internal/domain/entities"
)

// AuthMiddleware rejects requests that do not carry a valid Bearer JWT.
// On success it sets "userID", "username", and "role" in the Gin context.
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := extractClaims(c, jwtSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", string(claims.Role))
		c.Next()
	}
}

// OptionalAuthMiddleware parses the JWT when present but never blocks the
// request. Handlers can call c.Get("userID") to distinguish authenticated
// users from anonymous visitors (empty string == anonymous).
func OptionalAuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := extractClaims(c, jwtSecret)
		if err == nil {
			c.Set("userID", claims.UserID)
			c.Set("username", claims.Username)
			c.Set("role", string(claims.Role))
		} else {
			c.Set("userID", "")
			c.Set("username", "")
			c.Set("role", "")
		}
		c.Next()
	}
}

// AdminMiddleware must be placed after AuthMiddleware.
// It rejects requests from non-admin users with 403 Forbidden.
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role != string(entities.RoleAdmin) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "admin access required",
			})
			return
		}
		c.Next()
	}
}

// extractClaims parses and validates the Bearer token from the Authorization header.
func extractClaims(c *gin.Context, jwtSecret string) (*authusecase.Claims, error) {
	header := c.GetHeader("Authorization")
	if header == "" || !strings.HasPrefix(header, "Bearer ") {
		return nil, jwt.ErrTokenMalformed
	}

	tokenStr := strings.TrimPrefix(header, "Bearer ")

	claims := &authusecase.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, jwt.ErrTokenSignatureInvalid
	}

	return claims, nil
}
