package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"Smart_Task_Manager/internal/database" // Updated to use your wrapper
	"Smart_Task_Manager/internal/domain/entities"
	"Smart_Task_Manager/internal/repository/interfaces"
)

// userRepository is the concrete implementation of interfaces.UserRepository
type userRepository struct {
	db *database.DB // Changed from *sql.DB to *database.DB
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *database.DB) interfaces.UserRepository {
	return &userRepository{db: db}
}

// =========================================================================
// CRUD Operations
// =========================================================================

// Create saves a new user into the database
func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
	query := `
        INSERT INTO users (id, username, email, password_hash, created_at, updated_at) 
        VALUES ($1, $2, $3, $4, $5, $6)`

	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}
	user.UpdatedAt = user.CreatedAt

	// Using the DB wrapper which includes retry logic and metrics
	_, err := r.db.Exec(ctx, query,
		user.ID, user.Username, user.Email, user.PasswordHash,
		user.CreatedAt, user.UpdatedAt,
	)

	return err
}

// GetByID retrieves a user by their UUID
func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	query := `SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE id = $1`
	return r.fetchUser(ctx, query, id)
}

// GetByUsername retrieves a user by username (unique)
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*entities.User, error) {
	query := `SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE username = $1`
	return r.fetchUser(ctx, query, username)
}

// GetByEmail retrieves a user by email (unique)
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	query := `SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE email = $1`
	return r.fetchUser(ctx, query, email)
}

// Update modifies existing user details
func (r *userRepository) Update(ctx context.Context, user *entities.User) error {
	query := `
        UPDATE users SET 
            username = $1, email = $2, password_hash = $3, updated_at = $4
        WHERE id = $5`

	user.UpdatedAt = time.Now()

	// Using the DB wrapper
	tag, err := r.db.Exec(ctx, query,
		user.Username, user.Email, user.PasswordHash, user.UpdatedAt,
		user.ID,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errors.New("user not found")
	}

	return nil
}

// Delete removes a user completely from the database
func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errors.New("user not found")
	}

	return nil
}

// List retrieves all users
func (r *userRepository) List(ctx context.Context) ([]*entities.User, error) {
	query := `SELECT id, username, email, password_hash, created_at, updated_at FROM users ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entities.User
	for rows.Next() {
		user := &entities.User{}
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.PasswordHash,
			&user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

// =========================================================================
// Helper Methods
// =========================================================================

// fetchUser is a helper method to execute a query that returns a single user.
func (r *userRepository) fetchUser(ctx context.Context, query string, args ...any) (*entities.User, error) {
	user := &entities.User{}

	// Using the DB wrapper's QueryRow method
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}
