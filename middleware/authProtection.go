package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ProtectedRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				log.Println("invalid auth header")
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid auth header"})
				c.Abort()
				return
			}
			token = parts[1]
		} else {
			token = c.Query("token")
		}
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token Missing"})
			c.Abort()
			return
		}
		claims, err := ValidateToken(token)
		if err != nil {
			log.Println("validation error: ", err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}
		if claims.UserID == 0 {
			log.Println("user id absent")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user id absent in token"})
			c.Abort()
			return
		}
		log.Println("Token Validated")
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
