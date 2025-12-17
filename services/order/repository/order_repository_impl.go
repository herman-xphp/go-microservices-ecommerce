package repository

import (
	"github.com/herman-xphp/go-microservices-ecommerce/services/order/domain"
	"gorm.io/gorm"
)

type orderRepositoryImpl struct {
	db *gorm.DB
}

// NewOrderRepository creates a new instance of OrderRepository
func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepositoryImpl{db: db}
}

func (r *orderRepositoryImpl) Create(order *domain.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepositoryImpl) FindByID(id uint) (*domain.Order, error) {
	var order domain.Order
	err := r.db.Preload("Items").First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepositoryImpl) FindByUserID(userID uint, page, pageSize int) ([]domain.Order, int64, error) {
	var orders []domain.Order
	var total int64

	r.db.Model(&domain.Order{}).Where("user_id = ?", userID).Count(&total)

	offset := (page - 1) * pageSize
	err := r.db.Preload("Items").
		Where("user_id = ?", userID).
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&orders).Error

	return orders, total, err
}

func (r *orderRepositoryImpl) Update(order *domain.Order) error {
	return r.db.Save(order).Error
}

func (r *orderRepositoryImpl) UpdateStatus(id uint, status domain.OrderStatus) error {
	return r.db.Model(&domain.Order{}).Where("id = ?", id).Update("status", status).Error
}
