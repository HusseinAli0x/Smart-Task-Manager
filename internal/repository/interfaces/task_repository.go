package interfaces

import (
	"context"

	"github.com/google/uuid"

	"Smart_Task_Manager/internal/domain/entities"
)

// TaskRepository defines the database operations for a Task.
type TaskRepository interface {
	// Create saves a new task
	Create(ctx context.Context, task *entities.Task) error

	// GetByID retrieves a task by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Task, error)

	// Update modifies task details
	Update(ctx context.Context, task *entities.Task) error

	// Delete removes a task
	Delete(ctx context.Context, id uuid.UUID) error

	// ListByUserID retrieves all tasks for a specific user
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Task, error)

	// ListActiveTasks retrieves only incomplete tasks for a specific user (useful for dashboard)
	ListActiveTasks(ctx context.Context, userID uuid.UUID) ([]*entities.Task, error)
}
