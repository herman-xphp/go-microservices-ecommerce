package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/herman-xphp/go-microservices-ecommerce/pkg/utils"
	"github.com/herman-xphp/go-microservices-ecommerce/services/order/domain"
	"github.com/herman-xphp/go-microservices-ecommerce/services/order/dto"
	"github.com/herman-xphp/go-microservices-ecommerce/services/order/service"
)

type OrderHandler struct {
	orderService service.OrderService
}

// NewOrderHandler creates a new instance of OrderHandler
func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

// RegisterRoutes registers order routes to the gin router
func (h *OrderHandler) RegisterRoutes(router *gin.RouterGroup) {
	orders := router.Group("/orders")
	{
		orders.POST("", h.CreateOrder)
		orders.GET("", h.GetUserOrders)
		orders.GET("/:id", h.GetOrder)
		orders.PUT("/:id/status", h.UpdateOrderStatus)
		orders.POST("/:id/cancel", h.CancelOrder)
	}
}

// CreateOrder creates a new order
// POST /api/v1/orders
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	// Get user_id from context (set by auth middleware or from token)
	userIDVal, exists := c.Get("user_id")
	if !exists {
		// For demo purposes, allow passing user_id as query param
		userIDStr := c.Query("user_id")
		if userIDStr == "" {
			utils.ResponseError(c, http.StatusUnauthorized, "User ID required", nil)
			return
		}
		userID, _ := strconv.ParseUint(userIDStr, 10, 32)
		userIDVal = uint(userID)
	}

	userID := userIDVal.(uint)

	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	order, err := h.orderService.CreateOrder(c.Request.Context(), userID, &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrProductNotFound):
			utils.ResponseError(c, http.StatusBadRequest, "Product not found", err.Error())
		case errors.Is(err, service.ErrInsufficientStock):
			utils.ResponseError(c, http.StatusBadRequest, "Insufficient stock", err.Error())
		case errors.Is(err, service.ErrProductUnavailable):
			utils.ResponseError(c, http.StatusBadRequest, "Product unavailable", err.Error())
		default:
			utils.ResponseError(c, http.StatusInternalServerError, "Failed to create order", err.Error())
		}
		return
	}

	utils.ResponseSuccess(c, http.StatusCreated, "Order created successfully", order)
}

// GetOrder returns an order by ID
// GET /api/v1/orders/:id
func (h *OrderHandler) GetOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid order ID", nil)
		return
	}

	order, err := h.orderService.GetOrder(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, service.ErrOrderNotFound) {
			utils.ResponseError(c, http.StatusNotFound, "Order not found", nil)
			return
		}
		utils.ResponseError(c, http.StatusInternalServerError, "Failed to get order", err.Error())
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "Order retrieved successfully", order)
}

// GetUserOrders returns orders for the authenticated user
// GET /api/v1/orders?page=1&page_size=10
func (h *OrderHandler) GetUserOrders(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		userIDStr := c.Query("user_id")
		if userIDStr == "" {
			utils.ResponseError(c, http.StatusUnauthorized, "User ID required", nil)
			return
		}
		userID, _ := strconv.ParseUint(userIDStr, 10, 32)
		userIDVal = uint(userID)
	}

	userID := userIDVal.(uint)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	orders, err := h.orderService.GetUserOrders(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "Failed to get orders", err.Error())
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "Orders retrieved successfully", orders)
}

// UpdateOrderStatus updates the status of an order
// PUT /api/v1/orders/:id/status
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid order ID", nil)
		return
	}

	var req dto.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	err = h.orderService.UpdateOrderStatus(c.Request.Context(), uint(id), req.Status)
	if err != nil {
		if errors.Is(err, service.ErrOrderNotFound) {
			utils.ResponseError(c, http.StatusNotFound, "Order not found", nil)
			return
		}
		utils.ResponseError(c, http.StatusInternalServerError, "Failed to update order status", err.Error())
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "Order status updated successfully", gin.H{"status": req.Status})
}

// CancelOrder cancels a pending order
// POST /api/v1/orders/:id/cancel
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid order ID", nil)
		return
	}

	err = h.orderService.CancelOrder(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, service.ErrOrderNotFound) {
			utils.ResponseError(c, http.StatusNotFound, "Order not found", nil)
			return
		}
		utils.ResponseError(c, http.StatusBadRequest, "Failed to cancel order", err.Error())
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "Order cancelled successfully", gin.H{"status": domain.OrderStatusCancelled})
}
