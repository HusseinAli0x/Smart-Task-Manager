package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"Smart_Task_Manager/internal/utils"
)

// AuthMiddleware protects routes by ensuring a valid JWT is provided and tracks performance.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now() // Capture start time for performance tracking

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Authorization header is missing"})
			c.Abort()
			return
		}

		// Ensure the format is "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid format, use: Bearer <token>"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate the token and extract the UserID
		userID, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": err.Error()})
			c.Abort()
			return
		}

		// Inject the verified userID into the Gin context
		c.Set("userID", userID)

		// Log execution time for authentication check
		latency := time.Since(start)
		log.Printf("[METRIC] Auth Check | UserID: %v | Latency: %v", userID, latency)

		c.Next()
	}
}
