package entities

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"Smart_Task_Manager/internal/domain/enums"
)

type Task struct {
	ID          uuid.UUID           `db:"id" json:"id"`
	UserID      uuid.UUID           `db:"user_id" json:"user_id"`
	RawText     string              `db:"raw_text" json:"raw_text"`
	Title       string              `db:"title" json:"title"`
	Description *string             `db:"description" json:"description,omitempty"`
	Category    *enums.TaskCategory `db:"category" json:"category,omitempty"`
	Priority    enums.TaskPriority  `db:"priority" json:"priority"`
	DueDate     *time.Time          `db:"due_date" json:"due_date,omitempty"`
	SubTasks    json.RawMessage     `db:"sub_tasks" json:"sub_tasks"`
	IsCompleted bool                `db:"is_completed" json:"is_completed"`
	CreatedAt   time.Time           `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time           `db:"updated_at" json:"updated_at"`
}

// TableName returns the database table name for the entity
func (Task) TableName() string {
	return "tasks"
}

// IsOverdue checks if the task is past its due date and not yet completed
func (t *Task) IsOverdue() bool {
	if t.IsCompleted || t.DueDate == nil {
		return false
	}
	return t.DueDate.Before(time.Now())
}

// HasCategory checks if the task was assigned a category by the AI
func (t *Task) HasCategory() bool {
	return t.Category != nil
}

// HasDeadline checks if the task has a specific due date
func (t *Task) HasDeadline() bool {
	return t.DueDate != nil
}
