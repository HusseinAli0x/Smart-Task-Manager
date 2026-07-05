package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	"Smart_Task_Manager/internal/domain/entities"
	"Smart_Task_Manager/internal/repository/interfaces"
)

// taskRepository is the concrete implementation of interfaces.TaskRepository
type taskRepository struct {
	db *sql.DB
}

// NewTaskRepository creates a new instance of TaskRepository
func NewTaskRepository(db *sql.DB) interfaces.TaskRepository {
	return &taskRepository{db: db}
}

// =========================================================================
// CRUD Operations
// =========================================================================

// Create saves a new task into the database
func (r *taskRepository) Create(ctx context.Context, task *entities.Task) error {
	query := `
		INSERT INTO tasks (
			id, user_id, raw_text, title, description, 
			category, priority, due_date, sub_tasks, 
			is_completed, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, 
			$6, $7, $8, $9, 
			$10, $11, $12
		)`

	// Ensure timestamps are set before inserting
	if task.CreatedAt.IsZero() {
		task.CreatedAt = time.Now()
	}
	task.UpdatedAt = task.CreatedAt

	_, err := r.db.ExecContext(ctx, query,
		task.ID, task.UserID, task.RawText, task.Title, task.Description,
		task.Category, task.Priority, task.DueDate, task.SubTasks,
		task.IsCompleted, task.CreatedAt, task.UpdatedAt,
	)

	return err
}

// GetByID retrieves a task by its ID
func (r *taskRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Task, error) {
	query := `
		SELECT 
			id, user_id, raw_text, title, description, 
			category, priority, due_date, sub_tasks, 
			is_completed, created_at, updated_at
		FROM tasks 
		WHERE id = $1`

	task := &entities.Task{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&task.ID, &task.UserID, &task.RawText, &task.Title, &task.Description,
		&task.Category, &task.Priority, &task.DueDate, &task.SubTasks,
		&task.IsCompleted, &task.CreatedAt, &task.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("task not found") // يمكنك استبداله بـ ErrTaskNotFound الخاص بك
		}
		return nil, err
	}

	return task, nil
}

// Update modifies existing task details
func (r *taskRepository) Update(ctx context.Context, task *entities.Task) error {
	query := `
		UPDATE tasks SET 
			raw_text = $1, title = $2, description = $3, 
			category = $4, priority = $5, due_date = $6, 
			sub_tasks = $7, is_completed = $8, updated_at = $9
		WHERE id = $10`

	task.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		task.RawText, task.Title, task.Description,
		task.Category, task.Priority, task.DueDate,
		task.SubTasks, task.IsCompleted, task.UpdatedAt,
		task.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("task not found")
	}

	return nil
}

// Delete removes a task completely from the database
func (r *taskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM tasks WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("task not found")
	}

	return nil
}

// =========================================================================
// List Operations
// =========================================================================

// ListByUserID retrieves all tasks for a specific user
func (r *taskRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Task, error) {
	query := `
		SELECT 
			id, user_id, raw_text, title, description, 
			category, priority, due_date, sub_tasks, 
			is_completed, created_at, updated_at
		FROM tasks 
		WHERE user_id = $1 
		ORDER BY created_at DESC`

	return r.fetchTasks(ctx, query, userID)
}

// ListActiveTasks retrieves only incomplete tasks for a specific user
func (r *taskRepository) ListActiveTasks(ctx context.Context, userID uuid.UUID) ([]*entities.Task, error) {
	query := `
		SELECT 
			id, user_id, raw_text, title, description, 
			category, priority, due_date, sub_tasks, 
			is_completed, created_at, updated_at
		FROM tasks 
		WHERE user_id = $1 AND is_completed = FALSE 
		ORDER BY due_date ASC`

	return r.fetchTasks(ctx, query, userID)
}

// =========================================================================
// Helper Methods
// =========================================================================

// fetchTasks is a helper method to avoid repeating the rows iteration logic
func (r *taskRepository) fetchTasks(ctx context.Context, query string, args ...any) ([]*entities.Task, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*entities.Task
	for rows.Next() {
		task := &entities.Task{}
		err := rows.Scan(
			&task.ID, &task.UserID, &task.RawText, &task.Title, &task.Description,
			&task.Category, &task.Priority, &task.DueDate, &task.SubTasks,
			&task.IsCompleted, &task.CreatedAt, &task.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
