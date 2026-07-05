package handler

import "github.com/gin-gonic/gin"

// APIResponse defines the unified standard structure for all API outputs
type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

// respondWithJSON transforms data into a standard success JSON object
func respondWithJSON(c *gin.Context, code int, data any) {
	c.JSON(code, APIResponse{
		Success: true,
		Data:    data,
	})
}

// respondWithError transforms error messages into a standard failure JSON object
func respondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, APIResponse{
		Success: false,
		Message: message,
	})
}
