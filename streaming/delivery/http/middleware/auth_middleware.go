package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthMiddleware struct {
	jwtSecret string
}

type Claims struct {
	UserID uuid.UUID `json:"user_id"` // For backward compatibility
	Sub    string    `json:"sub"`     // Standard JWT claim (used by auth service)
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

func NewAuthMiddleware(jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: jwtSecret,
	}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(m.jwtSecret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
				"details": err.Error(),
			})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			// Get user ID from either 'sub' (standard) or 'user_id' (custom) claim
			var userID uuid.UUID
			var err error

			// Try to get from 'sub' first (used by auth service)
			if claims.Sub != "" {
				userID, err = uuid.Parse(claims.Sub)
				if err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{
						"error": "invalid token claims",
						"details": fmt.Sprintf("invalid user ID format in 'sub' claim: %v", err),
					})
					c.Abort()
					return
				}
			} else if claims.UserID != uuid.Nil {
				// Fallback to custom 'user_id' claim
				userID = claims.UserID
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "invalid token claims",
					"details": "token must contain either 'sub' or 'user_id' claim with a valid UUID",
				})
				c.Abort()
				return
			}

			// Validate user_id is not nil
			if userID == uuid.Nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "invalid token claims",
					"details": "token does not contain a valid user_id. Please ensure your authentication service generates tokens with user_id or sub claim.",
				})
				c.Abort()
				return
			}

			// Set user info in context
			c.Set("user_id", userID)
			c.Set("email", claims.Email)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			c.Abort()
			return
		}
	}
}

// Helper function to get user ID from context
func GetUserID(c *gin.Context) (uuid.UUID, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, fmt.Errorf("user_id not found in context")
	}

	id, ok := userID.(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid user_id type")
	}

	return id, nil
}
