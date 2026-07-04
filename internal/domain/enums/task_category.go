package enums

// TaskCategory represents standard AI-assigned task categories.
type TaskCategory string

const (
	CategoryWork     TaskCategory = "Work"
	CategoryPersonal TaskCategory = "Personal"
	CategoryStudy    TaskCategory = "Study"
)

// IsValid checks if the TaskCategory is allowed
func (c TaskCategory) IsValid() bool {
	switch c {
	case CategoryWork,
		CategoryPersonal,
		CategoryStudy:
		return true
	default:
		return false
	}
}
