package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"Smart_Task_Manager/internal/domain/entities"
	domainErrors "Smart_Task_Manager/internal/domain/errors"
	"Smart_Task_Manager/internal/repository/interfaces"
)

// =========================================================================
// Data Transfer Objects (DTOs)
// =========================================================================

// RegisterInput defines the required payloads to register a new user
type RegisterInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginInput defines the payloads required for user authentication
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserOutput defines the structured public data returned to the client
type UserOutput struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// LoginOutput defines the successful authentication response
type LoginOutput struct {
	User  *UserOutput `json:"user"`
	Token string      `json:"token"` // Placeholder for JWT or Session Token
}

// =========================================================================
// UseCase Interface
// =========================================================================

// UserUseCase coordinates the business logic operations for users
type UserUseCase interface {
	Register(ctx context.Context, input RegisterInput) (*UserOutput, error)
	Login(ctx context.Context, input LoginInput) (*LoginOutput, error)
}

// userUseCase is the concrete implementation of UserUseCase
type userUseCase struct {
	userRepo interfaces.UserRepository
}

// NewUserUseCase acts as a constructor to initialize userUseCase with its dependencies
func NewUserUseCase(repo interfaces.UserRepository) UserUseCase {
	return &userUseCase{userRepo: repo}
}

// =========================================================================
// Core Business Logic Implementation
// =========================================================================

// Register handles the business flow of creating a new secure user profile
func (u *userUseCase) Register(ctx context.Context, input RegisterInput) (*UserOutput, error) {
	// 1. Validate email uniqueness
	existingUser, err := u.userRepo.GetByEmail(ctx, input.Email)
	if err == nil && existingUser != nil {
		return nil, domainErrors.ErrEmailAlreadyExists
	}

	// 2. Validate username uniqueness
	existingUsername, err := u.userRepo.GetByUsername(ctx, input.Username)
	if err == nil && existingUsername != nil {
		return nil, domainErrors.ErrUsernameTaken
	}

	// 3. Encrypt the raw password securely using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 4. Transform DTO into domain entity
	user := &entities.User{
		ID:           uuid.New(),
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// 5. Persist entity data into the database layer
	err = u.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	// 6. Map entity to public output structure (sanitizing password hash)
	return &UserOutput{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

// Login handles user authentication logic and credentials verification
func (u *userUseCase) Login(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	// 1. Fetch user by email address
	user, err := u.userRepo.GetByEmail(ctx, input.Email)
	if err != nil || user == nil {
		return nil, domainErrors.ErrInvalidCredentials
	}

	// 2. Compare input password with the stored hash safely
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		return nil, domainErrors.ErrInvalidCredentials
	}

	// 3. Prepare the public user details output
	userPublicData := &UserOutput{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	// 4. Return user info alongside a token placeholder
	return &LoginOutput{
		User:  userPublicData,
		Token: "generated-jwt-token-placeholder", // Will be implemented with a real JWT service later
	}, nil
}
