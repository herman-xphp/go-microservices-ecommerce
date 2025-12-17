package repository

import (
	"github.com/herman-xphp/go-microservices-ecommerce/services/product/domain"
)

// MockProductRepository is a mock implementation for testing
type MockProductRepository struct {
	products map[uint]*domain.Product
	nextID   uint
}

func NewMockProductRepository() *MockProductRepository {
	return &MockProductRepository{
		products: make(map[uint]*domain.Product),
		nextID:   1,
	}
}

func (m *MockProductRepository) Create(product *domain.Product) error {
	product.ID = m.nextID
	m.nextID++
	m.products[product.ID] = product
	return nil
}

func (m *MockProductRepository) FindByID(id uint) (*domain.Product, error) {
	if product, ok := m.products[id]; ok {
		return product, nil
	}
	return nil, nil
}

func (m *MockProductRepository) FindAll(page, pageSize int) ([]domain.Product, int64, error) {
	var result []domain.Product
	for _, p := range m.products {
		if p.IsActive {
			result = append(result, *p)
		}
	}
	return result, int64(len(result)), nil
}

func (m *MockProductRepository) FindByCategory(categoryID uint, page, pageSize int) ([]domain.Product, int64, error) {
	var result []domain.Product
	for _, p := range m.products {
		if p.CategoryID == categoryID && p.IsActive {
			result = append(result, *p)
		}
	}
	return result, int64(len(result)), nil
}

func (m *MockProductRepository) Update(product *domain.Product) error {
	m.products[product.ID] = product
	return nil
}

func (m *MockProductRepository) Delete(id uint) error {
	delete(m.products, id)
	return nil
}

func (m *MockProductRepository) UpdateStock(id uint, quantity int) error {
	if product, ok := m.products[id]; ok {
		product.Stock += quantity
	}
	return nil
}

func (m *MockProductRepository) CheckStock(id uint) (int, error) {
	if product, ok := m.products[id]; ok {
		return product.Stock, nil
	}
	return 0, nil
}

// MockCategoryRepository is a mock implementation for testing
type MockCategoryRepository struct {
	categories map[uint]*domain.Category
	nextID     uint
}

func NewMockCategoryRepository() *MockCategoryRepository {
	return &MockCategoryRepository{
		categories: make(map[uint]*domain.Category),
		nextID:     1,
	}
}

func (m *MockCategoryRepository) Create(category *domain.Category) error {
	category.ID = m.nextID
	m.nextID++
	m.categories[category.ID] = category
	return nil
}

func (m *MockCategoryRepository) FindByID(id uint) (*domain.Category, error) {
	if cat, ok := m.categories[id]; ok {
		return cat, nil
	}
	return nil, nil
}

func (m *MockCategoryRepository) FindAll() ([]domain.Category, error) {
	var result []domain.Category
	for _, c := range m.categories {
		result = append(result, *c)
	}
	return result, nil
}

func (m *MockCategoryRepository) Update(category *domain.Category) error {
	m.categories[category.ID] = category
	return nil
}

func (m *MockCategoryRepository) Delete(id uint) error {
	delete(m.categories, id)
	return nil
}
