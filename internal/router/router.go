package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"Smart_Task_Manager/internal/handler"
	"Smart_Task_Manager/internal/middleware"
)

// Setup configures and organizes all application routes
func Setup(userHandler *handler.UserHandler, taskHandler *handler.TaskHandler) *gin.Engine {
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "Smart Task Manager API"})
	})

	// API route group
	api := r.Group("/api/v1")
	{
		// 1. Public authentication routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
		}

		// 2. Protected task routes (requires AuthMiddleware)
		tasks := api.Group("/tasks")
		tasks.Use(middleware.AuthMiddleware()) // Apply protection here
		{
			tasks.POST("", taskHandler.Create)
			tasks.GET("", taskHandler.ListUserTasks)
			tasks.GET("/active", taskHandler.ListActiveTasks)
			tasks.GET("/:id", taskHandler.GetByID)
			tasks.PUT("/:id", taskHandler.Update)
			tasks.DELETE("/:id", taskHandler.Delete)
		}
	}

	return r
}
