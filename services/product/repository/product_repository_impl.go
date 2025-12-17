package repository

import (
	"github.com/herman-xphp/go-microservices-ecommerce/services/product/domain"
	"gorm.io/gorm"
)

type productRepositoryImpl struct {
	db *gorm.DB
}

// NewProductRepository creates a new instance of ProductRepository
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepositoryImpl{db: db}
}

func (r *productRepositoryImpl) Create(product *domain.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepositoryImpl) FindByID(id uint) (*domain.Product, error) {
	var product domain.Product
	err := r.db.Preload("Category").First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepositoryImpl) FindAll(page, pageSize int) ([]domain.Product, int64, error) {
	var products []domain.Product
	var total int64

	r.db.Model(&domain.Product{}).Where("is_active = ?", true).Count(&total)

	offset := (page - 1) * pageSize
	err := r.db.Preload("Category").
		Where("is_active = ?", true).
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&products).Error

	return products, total, err
}

func (r *productRepositoryImpl) FindByCategory(categoryID uint, page, pageSize int) ([]domain.Product, int64, error) {
	var products []domain.Product
	var total int64

	r.db.Model(&domain.Product{}).
		Where("category_id = ? AND is_active = ?", categoryID, true).
		Count(&total)

	offset := (page - 1) * pageSize
	err := r.db.Preload("Category").
		Where("category_id = ? AND is_active = ?", categoryID, true).
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&products).Error

	return products, total, err
}

func (r *productRepositoryImpl) Update(product *domain.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&domain.Product{}, id).Error
}

func (r *productRepositoryImpl) UpdateStock(id uint, quantity int) error {
	return r.db.Model(&domain.Product{}).
		Where("id = ?", id).
		Update("stock", gorm.Expr("stock + ?", quantity)).Error
}

func (r *productRepositoryImpl) CheckStock(id uint) (int, error) {
	var product domain.Product
	err := r.db.Select("stock").First(&product, id).Error
	if err != nil {
		return 0, err
	}
	return product.Stock, nil
}

// Category Repository Implementation
type categoryRepositoryImpl struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new instance of CategoryRepository
func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepositoryImpl{db: db}
}

func (r *categoryRepositoryImpl) Create(category *domain.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepositoryImpl) FindByID(id uint) (*domain.Category, error) {
	var category domain.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepositoryImpl) FindAll() ([]domain.Category, error) {
	var categories []domain.Category
	err := r.db.Order("name ASC").Find(&categories).Error
	return categories, err
}

func (r *categoryRepositoryImpl) Update(category *domain.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&domain.Category{}, id).Error
}
