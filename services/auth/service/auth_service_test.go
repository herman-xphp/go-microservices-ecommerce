package service

import (
	"testing"

	"github.com/herman-xphp/go-microservices-ecommerce/services/auth/dto"
	"github.com/herman-xphp/go-microservices-ecommerce/services/auth/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testJWTSecret = "test-jwt-secret-for-unit-testing-min-32-chars"

func TestAuthService_Register_Success(t *testing.T) {
	// Arrange
	mockRepo := repository.NewMockUserRepository()
	authService := NewAuthService(mockRepo, testJWTSecret)

	req := &dto.RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	// Act
	resp, err := authService.Register(req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, "john@example.com", resp.User.Email)
	assert.Equal(t, "John Doe", resp.User.Name)
}

func TestAuthService_Register_DuplicateEmail(t *testing.T) {
	// Arrange
	mockRepo := repository.NewMockUserRepository()
	authService := NewAuthService(mockRepo, testJWTSecret)

	req := &dto.RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	// First registration
	_, err := authService.Register(req)
	require.NoError(t, err)

	// Act - second registration with same email
	resp, err := authService.Register(req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, ErrUserAlreadyExists, err)
}

func TestAuthService_Login_Success(t *testing.T) {
	// Arrange
	mockRepo := repository.NewMockUserRepository()
	authService := NewAuthService(mockRepo, testJWTSecret)

	// First register a user
	registerReq := &dto.RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}
	_, err := authService.Register(registerReq)
	require.NoError(t, err)

	// Act - login with correct credentials
	loginReq := &dto.LoginRequest{
		Email:    "john@example.com",
		Password: "password123",
	}
	resp, err := authService.Login(loginReq)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, "john@example.com", resp.User.Email)
}

func TestAuthService_Login_InvalidEmail(t *testing.T) {
	// Arrange
	mockRepo := repository.NewMockUserRepository()
	authService := NewAuthService(mockRepo, testJWTSecret)

	// Act - login with non-existent email
	loginReq := &dto.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}
	resp, err := authService.Login(loginReq)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, ErrInvalidCredentials, err)
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	// Arrange
	mockRepo := repository.NewMockUserRepository()
	authService := NewAuthService(mockRepo, testJWTSecret)

	// First register a user
	registerReq := &dto.RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}
	_, err := authService.Register(registerReq)
	require.NoError(t, err)

	// Act - login with wrong password
	loginReq := &dto.LoginRequest{
		Email:    "john@example.com",
		Password: "wrongpassword",
	}
	resp, err := authService.Login(loginReq)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, ErrInvalidCredentials, err)
}

func TestAuthService_ValidateToken_Success(t *testing.T) {
	// Arrange
	mockRepo := repository.NewMockUserRepository()
	authService := NewAuthService(mockRepo, testJWTSecret)

	// First register a user
	registerReq := &dto.RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}
	registerResp, err := authService.Register(registerReq)
	require.NoError(t, err)

	// Act - validate the token
	user, err := authService.ValidateToken(registerResp.Token)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "john@example.com", user.Email)
}

func TestAuthService_ValidateToken_InvalidToken(t *testing.T) {
	// Arrange
	mockRepo := repository.NewMockUserRepository()
	authService := NewAuthService(mockRepo, testJWTSecret)

	// Act - validate invalid token
	user, err := authService.ValidateToken("invalid-token")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestAuthService_GetUserByID_Success(t *testing.T) {
	// Arrange
	mockRepo := repository.NewMockUserRepository()
	authService := NewAuthService(mockRepo, testJWTSecret)

	// First register a user
	registerReq := &dto.RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}
	_, err := authService.Register(registerReq)
	require.NoError(t, err)

	// Act - get user by ID
	user, err := authService.GetUserByID(1)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "john@example.com", user.Email)
}

func TestAuthService_GetUserByID_NotFound(t *testing.T) {
	// Arrange
	mockRepo := repository.NewMockUserRepository()
	authService := NewAuthService(mockRepo, testJWTSecret)

	// Act - get non-existent user
	user, err := authService.GetUserByID(999)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
}
