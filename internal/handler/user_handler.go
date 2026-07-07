package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"Smart_Task_Manager/internal/usecase"
)

// UserHandler coordinates HTTP requests for user lifecycle operations using Gin
type UserHandler struct {
	userUseCase usecase.UserUseCase
}

// NewUserHandler initializes a new UserHandler instance with injected UseCase
func NewUserHandler(uuc usecase.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase: uuc}
}

// Register handles HTTP POST requests to create a new user profile
func (h *UserHandler) Register(c *gin.Context) {
	var input usecase.RegisterInput

	// Automatically bind incoming JSON payload into the input DTO struct
	if err := c.ShouldBindJSON(&input); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate basic input constraints
	if input.Email == "" || input.Username == "" || input.Password == "" {
		respondWithError(c, http.StatusBadRequest, "Missing required fields: username, email, password")
		return
	}

	// Trigger the registration business logic passing Gin's native context
	output, err := h.userUseCase.Register(c.Request.Context(), input)
	if err != nil {
		respondWithError(c, http.StatusConflict, err.Error())
		return
	}

	respondWithJSON(c, http.StatusCreated, output)
}

// Login handles HTTP POST requests for user authentication
func (h *UserHandler) Login(c *gin.Context) {
	var input usecase.LoginInput

	// Bind JSON data into the login DTO struct
	if err := c.ShouldBindJSON(&input); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.Email == "" || input.Password == "" {
		respondWithError(c, http.StatusBadRequest, "Missing email or password")
		return
	}

	// Execute authentication logic
	output, err := h.userUseCase.Login(c.Request.Context(), input)
	if err != nil {
		respondWithError(c, http.StatusUnauthorized, err.Error())
		return
	}

	respondWithJSON(c, http.StatusOK, output)
}
