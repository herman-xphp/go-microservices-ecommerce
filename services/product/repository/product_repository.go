package repository

import "github.com/herman-xphp/go-microservices-ecommerce/services/product/domain"

// ProductRepository defines the interface for product data operations
type ProductRepository interface {
	Create(product *domain.Product) error
	FindByID(id uint) (*domain.Product, error)
	FindAll(page, pageSize int) ([]domain.Product, int64, error)
	FindByCategory(categoryID uint, page, pageSize int) ([]domain.Product, int64, error)
	Update(product *domain.Product) error
	Delete(id uint) error
	UpdateStock(id uint, quantity int) error
	CheckStock(id uint) (int, error)
}

// CategoryRepository defines the interface for category data operations
type CategoryRepository interface {
	Create(category *domain.Category) error
	FindByID(id uint) (*domain.Category, error)
	FindAll() ([]domain.Category, error)
	Update(category *domain.Category) error
	Delete(id uint) error
}
