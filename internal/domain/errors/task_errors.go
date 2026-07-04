package errors

import "errors"

// Task-related errors
var (
	ErrTaskNotFound     = errors.New("task not found")
	ErrTaskUnauthorized = errors.New("user is not authorized to access or modify this task")
	ErrInvalidPriority  = errors.New("invalid task priority level")
	ErrInvalidCategory  = errors.New("invalid task category")
)
