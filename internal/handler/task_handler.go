package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"Smart_Task_Manager/internal/usecase"
)

// TaskHandler coordinates HTTP requests for task-related operations
type TaskHandler struct {
	taskUseCase usecase.TaskUseCase
}

// NewTaskHandler initializes a new TaskHandler instance with injected UseCase
func NewTaskHandler(tuc usecase.TaskUseCase) *TaskHandler {
	return &TaskHandler{taskUseCase: tuc}
}

// Create handles HTTP POST requests to create a new task
func (h *TaskHandler) Create(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		respondWithError(c, http.StatusUnauthorized, "Unauthorized: invalid or missing user context")
		return
	}

	var input usecase.CreateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Attach the authenticated user's ID to the input
	input.UserID = userID

	output, err := h.taskUseCase.Create(c.Request.Context(), input)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(c, http.StatusCreated, output)
}

// GetByID handles HTTP GET requests to retrieve a single task
func (h *TaskHandler) GetByID(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		respondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid task ID format")
		return
	}

	output, err := h.taskUseCase.GetByID(c.Request.Context(), taskID, userID)
	if err != nil {
		respondWithError(c, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(c, http.StatusOK, output)
}

// Update handles HTTP PUT requests to modify an existing task
func (h *TaskHandler) Update(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		respondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid task ID format")
		return
	}

	var input usecase.UpdateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Force correct IDs to prevent parameter tampering
	input.ID = taskID
	input.UserID = userID

	output, err := h.taskUseCase.Update(c.Request.Context(), input)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(c, http.StatusOK, output)
}

// Delete handles HTTP DELETE requests to remove a task
func (h *TaskHandler) Delete(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		respondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid task ID format")
		return
	}

	err = h.taskUseCase.Delete(c.Request.Context(), taskID, userID)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(c, http.StatusOK, map[string]string{"message": "Task deleted successfully"})
}

// ListUserTasks handles HTTP GET requests to retrieve all tasks for the authenticated user
func (h *TaskHandler) ListUserTasks(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		respondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	outputs, err := h.taskUseCase.ListUserTasks(c.Request.Context(), userID)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(c, http.StatusOK, outputs)
}

// ListActiveTasks handles HTTP GET requests to retrieve only active tasks
func (h *TaskHandler) ListActiveTasks(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		respondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	outputs, err := h.taskUseCase.ListActiveUserTasks(c.Request.Context(), userID)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(c, http.StatusOK, outputs)
}

// =========================================================================
// Helper Functions
// =========================================================================

// getUserIDFromContext safely extracts the UUID of the authenticated user from Gin's context
func getUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	val, exists := c.Get("userID")
	if !exists {
		return uuid.Nil, errors.New("userID not found in context")
	}

	userID, ok := val.(uuid.UUID)
	if !ok {
		// Attempt string parsing if it was stored as string
		idStr, isStr := val.(string)
		if !isStr {
			return uuid.Nil, errors.New("invalid userID type in context")
		}

		parsedID, err := uuid.Parse(idStr)
		if err != nil {
			return uuid.Nil, errors.New("malformed userID format")
		}
		return parsedID, nil
	}

	return userID, nil
}
