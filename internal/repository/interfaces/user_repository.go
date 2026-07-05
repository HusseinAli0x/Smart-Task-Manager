package interfaces

import (
	"context"

	"github.com/google/uuid"

	"Smart_Task_Manager/internal/domain/entities"
)

// UserRepository defines the database operations for a User.
type UserRepository interface {
	// Create saves a new user
	Create(ctx context.Context, user *entities.User) error

	// GetByID retrieves a user by their ID
	GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)

	// GetByUsername retrieves a user by username (unique)
	GetByUsername(ctx context.Context, username string) (*entities.User, error)

	// GetByEmail retrieves a user by email (unique)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)

	// Update modifies user details
	Update(ctx context.Context, user *entities.User) error

	// Delete removes a user
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves all users (simple list, no filters yet)
	List(ctx context.Context) ([]*entities.User, error)
}
