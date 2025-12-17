package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/herman-xphp/go-microservices-ecommerce/services/cart/dto"
	"github.com/herman-xphp/go-microservices-ecommerce/services/cart/service"
)

// CartHandler handles HTTP requests for cart operations
type CartHandler struct {
	cartService service.CartService
}

// NewCartHandler creates a new CartHandler
func NewCartHandler(cartService service.CartService) *CartHandler {
	return &CartHandler{cartService: cartService}
}

// RegisterRoutes registers cart routes
func (h *CartHandler) RegisterRoutes(router *gin.RouterGroup) {
	cart := router.Group("/cart")
	{
		cart.GET("", h.GetCart)
		cart.POST("/items", h.AddToCart)
		cart.PUT("/items/:product_id", h.UpdateItem)
		cart.DELETE("/items/:product_id", h.RemoveItem)
		cart.DELETE("", h.ClearCart)
	}
}

func (h *CartHandler) getUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		if uid := c.GetHeader("X-User-ID"); uid != "" {
			id, _ := strconv.ParseUint(uid, 10, 32)
			return uint(id), true
		}
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not authenticated",
		})
		return 0, false
	}
	return userID.(uint), true
}

// GetCart retrieves the user's cart
// GET /api/v1/cart
func (h *CartHandler) GetCart(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	cart, err := h.cartService.GetCart(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    cart,
	})
}

// AddToCart adds an item to the cart
// POST /api/v1/cart/items
func (h *CartHandler) AddToCart(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	var req dto.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	cart, err := h.cartService.AddToCart(c.Request.Context(), userID, &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err == service.ErrProductNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Item added to cart",
		"data":    cart,
	})
}

// UpdateItem updates item quantity in cart
// PUT /api/v1/cart/items/:product_id
func (h *CartHandler) UpdateItem(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	productID, err := strconv.ParseUint(c.Param("product_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid product ID",
		})
		return
	}

	var req dto.UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	cart, err := h.cartService.UpdateItem(c.Request.Context(), userID, uint(productID), &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err == service.ErrItemNotInCart {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Cart updated",
		"data":    cart,
	})
}

// RemoveItem removes an item from the cart
// DELETE /api/v1/cart/items/:product_id
func (h *CartHandler) RemoveItem(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	productID, err := strconv.ParseUint(c.Param("product_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid product ID",
		})
		return
	}

	cart, err := h.cartService.RemoveItem(c.Request.Context(), userID, uint(productID))
	if err != nil {
		status := http.StatusInternalServerError
		if err == service.ErrItemNotInCart {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Item removed from cart",
		"data":    cart,
	})
}

// ClearCart clears all items from the cart
// DELETE /api/v1/cart
func (h *CartHandler) ClearCart(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	if err := h.cartService.ClearCart(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Cart cleared",
	})
}
