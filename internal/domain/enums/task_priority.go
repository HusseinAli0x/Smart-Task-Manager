package enums

// TaskPriority represents the strict validation levels for a task's priority.
// This mirrors the PostgreSQL CHECK constraint (priority IN ('low', 'medium', 'high')).
type TaskPriority string

const (
	PriorityLow    TaskPriority = "low"
	PriorityMedium TaskPriority = "medium"
	PriorityHigh   TaskPriority = "high"
)

// IsValid checks if the TaskPriority is allowed
func (p TaskPriority) IsValid() bool {
	switch p {
	case PriorityLow,
		PriorityMedium,
		PriorityHigh:
		return true
	default:
		return false
	}
}
