package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/herman-xphp/go-microservices-ecommerce/pkg/utils"
	"github.com/herman-xphp/go-microservices-ecommerce/services/product/dto"
	"github.com/herman-xphp/go-microservices-ecommerce/services/product/service"
)

type ProductHandler struct {
	productService service.ProductService
}

// NewProductHandler creates a new instance of ProductHandler
func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

// RegisterRoutes registers product routes to the gin router
func (h *ProductHandler) RegisterRoutes(router *gin.RouterGroup) {
	products := router.Group("/products")
	{
		products.GET("", h.GetProducts)
		products.GET("/:id", h.GetProduct)
		products.POST("", h.CreateProduct)
		products.PUT("/:id", h.UpdateProduct)
		products.DELETE("/:id", h.DeleteProduct)
	}

	categories := router.Group("/categories")
	{
		categories.GET("", h.GetCategories)
		categories.POST("", h.CreateCategory)
	}
}

// GetProducts returns a paginated list of products
// GET /api/v1/products?page=1&page_size=10&category_id=1
func (h *ProductHandler) GetProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	categoryID, _ := strconv.Atoi(c.Query("category_id"))

	var response *dto.ProductListResponse
	var err error

	if categoryID > 0 {
		response, err = h.productService.GetProductsByCategory(uint(categoryID), page, pageSize)
	} else {
		response, err = h.productService.GetProducts(page, pageSize)
	}

	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "Failed to get products", err.Error())
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "Products retrieved successfully", response)
}

// GetProduct returns a single product by ID
// GET /api/v1/products/:id
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid product ID", nil)
		return
	}

	product, err := h.productService.GetProduct(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			utils.ResponseError(c, http.StatusNotFound, "Product not found", nil)
			return
		}
		utils.ResponseError(c, http.StatusInternalServerError, "Failed to get product", err.Error())
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "Product retrieved successfully", product)
}

// CreateProduct creates a new product
// POST /api/v1/products
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	product, err := h.productService.CreateProduct(&req)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "Failed to create product", err.Error())
		return
	}

	utils.ResponseSuccess(c, http.StatusCreated, "Product created successfully", product)
}

// UpdateProduct updates an existing product
// PUT /api/v1/products/:id
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid product ID", nil)
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	product, err := h.productService.UpdateProduct(uint(id), &req)
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			utils.ResponseError(c, http.StatusNotFound, "Product not found", nil)
			return
		}
		utils.ResponseError(c, http.StatusInternalServerError, "Failed to update product", err.Error())
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "Product updated successfully", product)
}

// DeleteProduct deletes a product
// DELETE /api/v1/products/:id
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid product ID", nil)
		return
	}

	err = h.productService.DeleteProduct(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			utils.ResponseError(c, http.StatusNotFound, "Product not found", nil)
			return
		}
		utils.ResponseError(c, http.StatusInternalServerError, "Failed to delete product", err.Error())
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "Product deleted successfully", nil)
}

// GetCategories returns all categories
// GET /api/v1/categories
func (h *ProductHandler) GetCategories(c *gin.Context) {
	categories, err := h.productService.GetCategories()
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "Failed to get categories", err.Error())
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "Categories retrieved successfully", categories)
}

// CreateCategory creates a new category
// POST /api/v1/categories
func (h *ProductHandler) CreateCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	category, err := h.productService.CreateCategory(&req)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "Failed to create category", err.Error())
		return
	}

	utils.ResponseSuccess(c, http.StatusCreated, "Category created successfully", category)
}
