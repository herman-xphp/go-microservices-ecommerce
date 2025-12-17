package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/herman-xphp/go-microservices-ecommerce/services/gateway/client"
)

// GatewayHandler provides aggregation endpoints
type GatewayHandler struct {
	authClient    *client.AuthClient
	productClient *client.ProductClient
}

// NewGatewayHandler creates a new gateway handler
func NewGatewayHandler(authClient *client.AuthClient, productClient *client.ProductClient) *GatewayHandler {
	return &GatewayHandler{
		authClient:    authClient,
		productClient: productClient,
	}
}

// GetUserProfile returns user profile with additional data
// GET /api/v1/me
func (h *GatewayHandler) GetUserProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not authenticated",
		})
		return
	}

	user, err := h.authClient.GetUserByID(c.Request.Context(), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get user profile",
			"error":   err.Error(),
		})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile retrieved successfully",
		"data": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	})
}

// GetProductWithStock returns product with real-time stock info
// GET /api/v1/products/:id/stock
func (h *GatewayHandler) GetProductWithStock(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid product ID",
		})
		return
	}

	// Get product info via gRPC
	product, err := h.productClient.GetProduct(c.Request.Context(), uint(productID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get product",
			"error":   err.Error(),
		})
		return
	}

	if product == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Product not found",
		})
		return
	}

	// Get real-time stock via gRPC
	stock, found, err := h.productClient.CheckStock(c.Request.Context(), uint(productID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to check stock",
			"error":   err.Error(),
		})
		return
	}

	if !found {
		stock = 0
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Product retrieved with stock info",
		"data": gin.H{
			"id":          product.ID,
			"name":        product.Name,
			"description": product.Description,
			"price":       product.Price,
			"stock":       stock,
			"is_active":   product.IsActive,
			"in_stock":    stock > 0,
		},
	})
}
