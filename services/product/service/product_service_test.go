package service

import (
	"testing"

	"github.com/herman-xphp/go-microservices-ecommerce/services/product/dto"
	"github.com/herman-xphp/go-microservices-ecommerce/services/product/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductService_CreateProduct_Success(t *testing.T) {
	// Arrange
	productRepo := repository.NewMockProductRepository()
	categoryRepo := repository.NewMockCategoryRepository()
	productService := NewProductService(productRepo, categoryRepo)

	req := &dto.CreateProductRequest{
		Name:        "Test Product",
		Description: "A test product",
		Price:       99.99,
		Stock:       100,
		CategoryID:  1,
	}

	// Act
	resp, err := productService.CreateProduct(req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Test Product", resp.Name)
	assert.Equal(t, 99.99, resp.Price)
	assert.Equal(t, 100, resp.Stock)
	assert.True(t, resp.IsActive)
}

func TestProductService_GetProduct_Success(t *testing.T) {
	// Arrange
	productRepo := repository.NewMockProductRepository()
	categoryRepo := repository.NewMockCategoryRepository()
	productService := NewProductService(productRepo, categoryRepo)

	// Create a product first
	createReq := &dto.CreateProductRequest{
		Name:  "Test Product",
		Price: 99.99,
		Stock: 100,
	}
	created, err := productService.CreateProduct(createReq)
	require.NoError(t, err)

	// Act
	resp, err := productService.GetProduct(created.ID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, created.ID, resp.ID)
	assert.Equal(t, "Test Product", resp.Name)
}

func TestProductService_GetProduct_NotFound(t *testing.T) {
	// Arrange
	productRepo := repository.NewMockProductRepository()
	categoryRepo := repository.NewMockCategoryRepository()
	productService := NewProductService(productRepo, categoryRepo)

	// Act
	resp, err := productService.GetProduct(999)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, ErrProductNotFound, err)
}

func TestProductService_CheckStock_Success(t *testing.T) {
	// Arrange
	productRepo := repository.NewMockProductRepository()
	categoryRepo := repository.NewMockCategoryRepository()
	productService := NewProductService(productRepo, categoryRepo)

	// Create a product
	createReq := &dto.CreateProductRequest{
		Name:  "Test Product",
		Price: 99.99,
		Stock: 50,
	}
	created, err := productService.CreateProduct(createReq)
	require.NoError(t, err)

	// Act
	stock, err := productService.CheckStock(created.ID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 50, stock)
}

func TestProductService_DecreaseStock_Success(t *testing.T) {
	// Arrange
	productRepo := repository.NewMockProductRepository()
	categoryRepo := repository.NewMockCategoryRepository()
	productService := NewProductService(productRepo, categoryRepo)

	// Create a product
	createReq := &dto.CreateProductRequest{
		Name:  "Test Product",
		Price: 99.99,
		Stock: 50,
	}
	created, err := productService.CreateProduct(createReq)
	require.NoError(t, err)

	// Act
	err = productService.DecreaseStock(created.ID, 10)

	// Assert
	require.NoError(t, err)

	// Verify stock decreased
	stock, err := productService.CheckStock(created.ID)
	require.NoError(t, err)
	assert.Equal(t, 40, stock)
}

func TestProductService_DecreaseStock_InsufficientStock(t *testing.T) {
	// Arrange
	productRepo := repository.NewMockProductRepository()
	categoryRepo := repository.NewMockCategoryRepository()
	productService := NewProductService(productRepo, categoryRepo)

	// Create a product
	createReq := &dto.CreateProductRequest{
		Name:  "Test Product",
		Price: 99.99,
		Stock: 10,
	}
	created, err := productService.CreateProduct(createReq)
	require.NoError(t, err)

	// Act - try to decrease more than available
	err = productService.DecreaseStock(created.ID, 20)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrInsufficientStock, err)
}

func TestProductService_CreateCategory_Success(t *testing.T) {
	// Arrange
	productRepo := repository.NewMockProductRepository()
	categoryRepo := repository.NewMockCategoryRepository()
	productService := NewProductService(productRepo, categoryRepo)

	req := &dto.CreateCategoryRequest{
		Name: "Electronics",
	}

	// Act
	resp, err := productService.CreateCategory(req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Electronics", resp.Name)
}
