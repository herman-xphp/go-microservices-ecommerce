package repository

import "github.com/herman-xphp/go-microservices-ecommerce/services/payment/domain"

// PaymentRepository defines the interface for payment data operations
type PaymentRepository interface {
	Create(payment *domain.Payment) error
	FindByID(id uint) (*domain.Payment, error)
	FindByOrderID(orderID uint) (*domain.Payment, error)
	FindByTransactionID(transactionID string) (*domain.Payment, error)
	FindByUserID(userID uint, page, pageSize int) ([]domain.Payment, int64, error)
	Update(payment *domain.Payment) error
	UpdateStatus(id uint, status domain.PaymentStatus) error
}
