package rest

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const userKey = "user"

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Code:    "missing-token",
				Message: "Brak tokenu autoryzacyjnego",
			})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("nieoczekiwana metoda podpisu: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Code:    "invalid-token",
				Message: "Nieprawidłowy lub wygasły token",
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Code:    "invalid-token-claims",
				Message: "Nieprawidłowy format tokenu",
			})
			return
		}

		userID, err := uuid.Parse(claims["sub"].(string))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Code:    "invalid-user-id",
				Message: "Nieprawidłowy identyfikator użytkownika w tokenie",
			})
			return
		}

		ctx := context.WithValue(c.Request.Context(), userKey, userID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(userKey).(uuid.UUID)
	return userID, ok
}
