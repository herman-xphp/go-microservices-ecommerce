package repository

import (
	"github.com/herman-xphp/go-microservices-ecommerce/services/payment/domain"
	"gorm.io/gorm"
)

type paymentRepositoryImpl struct {
	db *gorm.DB
}

// NewPaymentRepository creates a new instance of PaymentRepository
func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepositoryImpl{db: db}
}

func (r *paymentRepositoryImpl) Create(payment *domain.Payment) error {
	return r.db.Create(payment).Error
}

func (r *paymentRepositoryImpl) FindByID(id uint) (*domain.Payment, error) {
	var payment domain.Payment
	err := r.db.First(&payment, id).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepositoryImpl) FindByOrderID(orderID uint) (*domain.Payment, error) {
	var payment domain.Payment
	err := r.db.Where("order_id = ?", orderID).First(&payment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepositoryImpl) FindByTransactionID(transactionID string) (*domain.Payment, error) {
	var payment domain.Payment
	err := r.db.Where("transaction_id = ?", transactionID).First(&payment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepositoryImpl) FindByUserID(userID uint, page, pageSize int) ([]domain.Payment, int64, error) {
	var payments []domain.Payment
	var total int64

	r.db.Model(&domain.Payment{}).Where("user_id = ?", userID).Count(&total)

	offset := (page - 1) * pageSize
	err := r.db.Where("user_id = ?", userID).
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&payments).Error

	return payments, total, err
}

func (r *paymentRepositoryImpl) Update(payment *domain.Payment) error {
	return r.db.Save(payment).Error
}

func (r *paymentRepositoryImpl) UpdateStatus(id uint, status domain.PaymentStatus) error {
	return r.db.Model(&domain.Payment{}).Where("id = ?", id).Update("status", status).Error
}
