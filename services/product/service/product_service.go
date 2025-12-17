package service

import (
	"errors"
	"math"

	"github.com/herman-xphp/go-microservices-ecommerce/services/product/domain"
	"github.com/herman-xphp/go-microservices-ecommerce/services/product/dto"
	"github.com/herman-xphp/go-microservices-ecommerce/services/product/repository"
	"gorm.io/gorm"
)

var (
	ErrProductNotFound   = errors.New("product not found")
	ErrCategoryNotFound  = errors.New("category not found")
	ErrInsufficientStock = errors.New("insufficient stock")
)

// ProductService defines the interface for product operations
type ProductService interface {
	// Product CRUD
	CreateProduct(req *dto.CreateProductRequest) (*dto.ProductResponse, error)
	GetProduct(id uint) (*dto.ProductResponse, error)
	GetProducts(page, pageSize int) (*dto.ProductListResponse, error)
	GetProductsByCategory(categoryID uint, page, pageSize int) (*dto.ProductListResponse, error)
	UpdateProduct(id uint, req *dto.UpdateProductRequest) (*dto.ProductResponse, error)
	DeleteProduct(id uint) error

	// Stock operations (for gRPC)
	CheckStock(productID uint) (int, error)
	DecreaseStock(productID uint, quantity int) error

	// Category CRUD
	CreateCategory(req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error)
	GetCategories() ([]dto.CategoryResponse, error)
}

type productServiceImpl struct {
	productRepo  repository.ProductRepository
	categoryRepo repository.CategoryRepository
}

// NewProductService creates a new instance of ProductService
func NewProductService(productRepo repository.ProductRepository, categoryRepo repository.CategoryRepository) ProductService {
	return &productServiceImpl{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *productServiceImpl) CreateProduct(req *dto.CreateProductRequest) (*dto.ProductResponse, error) {
	product := &domain.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		CategoryID:  req.CategoryID,
		ImageURL:    req.ImageURL,
		IsActive:    true,
	}

	if err := s.productRepo.Create(product); err != nil {
		return nil, err
	}

	return s.toProductResponse(product), nil
}

func (s *productServiceImpl) GetProduct(id uint) (*dto.ProductResponse, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	if product == nil {
		return nil, ErrProductNotFound
	}
	return s.toProductResponse(product), nil
}

func (s *productServiceImpl) GetProducts(page, pageSize int) (*dto.ProductListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	products, total, err := s.productRepo.FindAll(page, pageSize)
	if err != nil {
		return nil, err
	}

	return s.toProductListResponse(products, total, page, pageSize), nil
}

func (s *productServiceImpl) GetProductsByCategory(categoryID uint, page, pageSize int) (*dto.ProductListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	products, total, err := s.productRepo.FindByCategory(categoryID, page, pageSize)
	if err != nil {
		return nil, err
	}

	return s.toProductListResponse(products, total, page, pageSize), nil
}

func (s *productServiceImpl) UpdateProduct(id uint, req *dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}
	if req.CategoryID != nil {
		product.CategoryID = *req.CategoryID
	}
	if req.ImageURL != nil {
		product.ImageURL = *req.ImageURL
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	if err := s.productRepo.Update(product); err != nil {
		return nil, err
	}

	return s.toProductResponse(product), nil
}

func (s *productServiceImpl) DeleteProduct(id uint) error {
	_, err := s.productRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProductNotFound
		}
		return err
	}
	return s.productRepo.Delete(id)
}

func (s *productServiceImpl) CheckStock(productID uint) (int, error) {
	stock, err := s.productRepo.CheckStock(productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrProductNotFound
		}
		return 0, err
	}
	return stock, nil
}

func (s *productServiceImpl) DecreaseStock(productID uint, quantity int) error {
	stock, err := s.CheckStock(productID)
	if err != nil {
		return err
	}

	if stock < quantity {
		return ErrInsufficientStock
	}

	return s.productRepo.UpdateStock(productID, -quantity)
}

func (s *productServiceImpl) CreateCategory(req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	category := &domain.Category{
		Name: req.Name,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
	}, nil
}

func (s *productServiceImpl) GetCategories() ([]dto.CategoryResponse, error) {
	categories, err := s.categoryRepo.FindAll()
	if err != nil {
		return nil, err
	}

	result := make([]dto.CategoryResponse, len(categories))
	for i, cat := range categories {
		result[i] = dto.CategoryResponse{
			ID:   cat.ID,
			Name: cat.Name,
		}
	}
	return result, nil
}

// Helper methods
func (s *productServiceImpl) toProductResponse(p *domain.Product) *dto.ProductResponse {
	resp := &dto.ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		CategoryID:  p.CategoryID,
		ImageURL:    p.ImageURL,
		IsActive:    p.IsActive,
	}

	if p.Category != nil {
		resp.Category = &dto.CategoryResponse{
			ID:   p.Category.ID,
			Name: p.Category.Name,
		}
	}

	return resp
}

func (s *productServiceImpl) toProductListResponse(products []domain.Product, total int64, page, pageSize int) *dto.ProductListResponse {
	productResponses := make([]dto.ProductResponse, len(products))
	for i, p := range products {
		productResponses[i] = *s.toProductResponse(&p)
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.ProductListResponse{
		Products:   productResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}
