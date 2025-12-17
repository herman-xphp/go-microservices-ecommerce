package repository

import "github.com/herman-xphp/go-microservices-ecommerce/services/order/domain"

// OrderRepository defines the interface for order data operations
type OrderRepository interface {
	Create(order *domain.Order) error
	FindByID(id uint) (*domain.Order, error)
	FindByUserID(userID uint, page, pageSize int) ([]domain.Order, int64, error)
	Update(order *domain.Order) error
	UpdateStatus(id uint, status domain.OrderStatus) error
}
