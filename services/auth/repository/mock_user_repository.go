package repository

import (
	"github.com/herman-xphp/go-microservices-ecommerce/services/auth/domain"
)

// MockUserRepository is a mock implementation of UserRepository for testing
type MockUserRepository struct {
	users map[string]*domain.User
	byID  map[uint]*domain.User
}

// NewMockUserRepository creates a new mock user repository
func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*domain.User),
		byID:  make(map[uint]*domain.User),
	}
}

func (m *MockUserRepository) Create(user *domain.User) error {
	user.ID = uint(len(m.users) + 1)
	m.users[user.Email] = user
	m.byID[user.ID] = user
	return nil
}

func (m *MockUserRepository) FindByEmail(email string) (*domain.User, error) {
	if user, ok := m.users[email]; ok {
		return user, nil
	}
	return nil, nil
}

func (m *MockUserRepository) FindByID(id uint) (*domain.User, error) {
	if user, ok := m.byID[id]; ok {
		return user, nil
	}
	return nil, nil
}

// AddUser adds a user directly to the mock (for testing)
func (m *MockUserRepository) AddUser(user *domain.User) {
	m.users[user.Email] = user
	m.byID[user.ID] = user
}
