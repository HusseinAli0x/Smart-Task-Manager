package usecase

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"Smart_Task_Manager/internal/domain/entities"
	"Smart_Task_Manager/internal/domain/enums"
	domainErrors "Smart_Task_Manager/internal/domain/errors"
	"Smart_Task_Manager/internal/repository/interfaces"
)

// =========================================================================
// Data Transfer Objects (DTOs)
// =========================================================================

// CreateTaskInput defines the payload required to create a new smart task
type CreateTaskInput struct {
	UserID      uuid.UUID       `json:"user_id"`
	RawText     string          `json:"raw_text"`
	Title       string          `json:"title"`
	Description *string         `json:"description,omitempty"`
	Category    *string         `json:"category,omitempty"`
	Priority    string          `json:"priority"` // Expected: "low", "medium", "high"
	DueDate     *time.Time      `json:"due_date,omitempty"`
	SubTasks    json.RawMessage `json:"sub_tasks"`
}

// UpdateTaskInput defines the payload required to modify an existing task
type UpdateTaskInput struct {
	ID          uuid.UUID       `json:"id"`
	UserID      uuid.UUID       `json:"user_id"` // Used to verify ownership before update
	Title       string          `json:"title"`
	Description *string         `json:"description,omitempty"`
	Category    *string         `json:"category,omitempty"`
	Priority    string          `json:"priority"`
	DueDate     *time.Time      `json:"due_date,omitempty"`
	SubTasks    json.RawMessage `json:"sub_tasks"`
	IsCompleted bool            `json:"is_completed"`
}

// TaskOutput defines the sanitized structured response for client consumption
type TaskOutput struct {
	ID          uuid.UUID           `json:"id"`
	UserID      uuid.UUID           `json:"user_id"`
	RawText     string              `json:"raw_text"`
	Title       string              `json:"title"`
	Description *string             `json:"description,omitempty"`
	Category    *enums.TaskCategory `json:"category,omitempty"`
	Priority    enums.TaskPriority  `json:"priority"`
	DueDate     *time.Time          `json:"due_date,omitempty"`
	SubTasks    json.RawMessage     `json:"sub_tasks"`
	IsCompleted bool                `json:"is_completed"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

// =========================================================================
// UseCase Interface
// =========================================================================

// TaskUseCase coordinates the core business operations for smart tasks
type TaskUseCase interface {
	Create(ctx context.Context, input CreateTaskInput) (*TaskOutput, error)
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*TaskOutput, error)
	Update(ctx context.Context, input UpdateTaskInput) (*TaskOutput, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	ListUserTasks(ctx context.Context, userID uuid.UUID) ([]*TaskOutput, error)
	ListActiveUserTasks(ctx context.Context, userID uuid.UUID) ([]*TaskOutput, error)
}

// taskUseCase is the concrete implementation of TaskUseCase
type taskUseCase struct {
	taskRepo interfaces.TaskRepository
}

// NewTaskUseCase initializes a new taskUseCase with injected dependencies
func NewTaskUseCase(repo interfaces.TaskRepository) TaskUseCase {
	return &taskUseCase{taskRepo: repo}
}

// =========================================================================
// Core Business Logic Implementation
// =========================================================================

// Create validates and stores a newly extracted AI smart task
func (u *taskUseCase) Create(ctx context.Context, input CreateTaskInput) (*TaskOutput, error) {
	// 1. Convert and validate priority enum value
	taskPriority := enums.TaskPriority(input.Priority)
	if !taskPriority.IsValid() {
		return nil, domainErrors.ErrInvalidPriority
	}

	// 2. Safely map optional category if provided
	var taskCategory *enums.TaskCategory
	if input.Category != nil {
		cat := enums.TaskCategory(*input.Category)
		if !cat.IsValid() {
			return nil, domainErrors.ErrInvalidCategory
		}
		taskCategory = &cat
	}

	// 3. Fallback to empty array if subtasks are nil
	subTasks := input.SubTasks
	if subTasks == nil {
		subTasks = json.RawMessage("[]")
	}

	// 4. Transform DTO payload into domain entity
	task := &entities.Task{
		ID:          uuid.New(),
		UserID:      input.UserID,
		RawText:     input.RawText,
		Title:       input.Title,
		Description: input.Description,
		Category:    taskCategory,
		Priority:    taskPriority,
		DueDate:     input.DueDate,
		SubTasks:    subTasks,
		IsCompleted: false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 5. Save the task using repository layer
	err := u.taskRepo.Create(ctx, task)
	if err != nil {
		return nil, err
	}

	return u.mapToOutput(task), nil
}

// GetByID fetches a single task, enforcing strict user ownership validation
func (u *taskUseCase) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*TaskOutput, error) {
	task, err := u.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domainErrors.ErrTaskNotFound
	}

	// Security Check: Enforce data isolation between separate users
	if task.UserID != userID {
		return nil, domainErrors.ErrTaskUnauthorized
	}

	return u.mapToOutput(task), nil
}

// Update handles editing task info and modification permissions mapping
func (u *taskUseCase) Update(ctx context.Context, input UpdateTaskInput) (*TaskOutput, error) {
	// 1. Fetch existing task to ensure presence and check authorization
	task, err := u.taskRepo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, domainErrors.ErrTaskNotFound
	}

	if task.UserID != input.UserID {
		return nil, domainErrors.ErrTaskUnauthorized
	}

	// 2. Validate incoming enum modifications
	taskPriority := enums.TaskPriority(input.Priority)
	if !taskPriority.IsValid() {
		return nil, domainErrors.ErrInvalidPriority
	}

	var taskCategory *enums.TaskCategory
	if input.Category != nil {
		cat := enums.TaskCategory(*input.Category)
		if !cat.IsValid() {
			return nil, domainErrors.ErrInvalidCategory
		}
		taskCategory = &cat
	}

	// 3. Mutate domain entity state
	task.Title = input.Title
	task.Description = input.Description
	task.Category = taskCategory
	task.Priority = taskPriority
	task.DueDate = input.DueDate
	task.IsCompleted = input.IsCompleted
	if input.SubTasks != nil {
		task.SubTasks = input.SubTasks
	}

	// 4. Update records inside the repository
	err = u.taskRepo.Update(ctx, task)
	if err != nil {
		return nil, err
	}

	return u.mapToOutput(task), nil
}

// Delete permanently purges a specific task from storage
func (u *taskUseCase) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	task, err := u.taskRepo.GetByID(ctx, id)
	if err != nil {
		return domainErrors.ErrTaskNotFound
	}

	if task.UserID != userID {
		return domainErrors.ErrTaskUnauthorized
	}

	return u.taskRepo.Delete(ctx, id)
}

// ListUserTasks extracts chronological history records for a distinct profile
func (u *taskUseCase) ListUserTasks(ctx context.Context, userID uuid.UUID) ([]*TaskOutput, error) {
	tasks, err := u.taskRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return u.mapListToOutput(tasks), nil
}

// ListActiveUserTasks extracts incomplete deadline-oriented tasks for dashboards
func (u *taskUseCase) ListActiveUserTasks(ctx context.Context, userID uuid.UUID) ([]*TaskOutput, error) {
	tasks, err := u.taskRepo.ListActiveTasks(ctx, userID) // Note: fixed variable reference from 'r' to 'u'
	if err != nil {
		return nil, err
	}

	return u.mapListToOutput(tasks), nil
}

// =========================================================================
// Private Translation Helpers
// =========================================================================

func (u *taskUseCase) mapToOutput(task *entities.Task) *TaskOutput {
	return &TaskOutput{
		ID:          task.ID,
		UserID:      task.UserID,
		RawText:     task.RawText,
		Title:       task.Title,
		Description: task.Description,
		Category:    task.Category,
		Priority:    task.Priority,
		DueDate:     task.DueDate,
		SubTasks:    task.SubTasks,
		IsCompleted: task.IsCompleted,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}

func (u *taskUseCase) mapListToOutput(tasks []*entities.Task) []*TaskOutput {
	outputs := make([]*TaskOutput, 0, len(tasks))
	for _, task := range tasks {
		outputs = append(outputs, u.mapToOutput(task))
	}
	return outputs
}
