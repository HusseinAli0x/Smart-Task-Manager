package main

import (
	"context"
	"log"
	"time" // Import time package to track startup duration

	"Smart_Task_Manager/internal/config"
	"Smart_Task_Manager/internal/database"
	"Smart_Task_Manager/internal/handler"
	"Smart_Task_Manager/internal/repository/postgres"
	"Smart_Task_Manager/internal/router"
	"Smart_Task_Manager/internal/usecase"
)

func main() {
	// Capture the start time of the application
	startTime := time.Now()

	// 1. Load application configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Could not load configuration: %v", err)
	}

	// 2. Initialize database wrapper
	dbInstance, err := database.New(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbInstance.Close()

	// Verify database connectivity
	if err := dbInstance.Ping(context.Background()); err != nil {
		log.Fatalf("Database is unreachable: %v", err)
	}
	log.Println("Database connection established successfully!")

	// 3. Initialize layers (Repositories, UseCases, Handlers)
	userRepo := postgres.NewUserRepository(dbInstance)
	taskRepo := postgres.NewTaskRepository(dbInstance)

	userUseCase := usecase.NewUserUseCase(userRepo)
	taskUseCase := usecase.NewTaskUseCase(taskRepo)

	userHandler := handler.NewUserHandler(userUseCase)
	taskHandler := handler.NewTaskHandler(taskUseCase)

	// 4. Setup router
	r := router.Setup(userHandler, taskHandler)

	// Calculate total startup latency
	startupLatency := time.Since(startTime)
	log.Printf("Application initialized in: %v", startupLatency)

	// 5. Start the HTTP server
	serverAddr := cfg.Server.Address()
	log.Printf("Smart Task Manager API starting on %s...", serverAddr)

	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
