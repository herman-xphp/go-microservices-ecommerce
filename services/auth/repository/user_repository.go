package repository

import "github.com/username/go-microservices-ecommerce/services/auth/domain"

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
	FindByID(id uint) (*domain.User, error)
}
