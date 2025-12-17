package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/herman-xphp/go-microservices-ecommerce/services/payment/dto"
	"github.com/herman-xphp/go-microservices-ecommerce/services/payment/service"
)

// PaymentHandler handles HTTP requests for payment operations
type PaymentHandler struct {
	paymentService service.PaymentService
}

// NewPaymentHandler creates a new PaymentHandler
func NewPaymentHandler(paymentService service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

// RegisterRoutes registers payment-related routes
func (h *PaymentHandler) RegisterRoutes(router *gin.RouterGroup) {
	payments := router.Group("/payments")
	{
		payments.POST("", h.CreatePayment)
		payments.GET("", h.GetUserPayments)
		payments.GET("/:id", h.GetPayment)
		payments.GET("/order/:order_id", h.GetPaymentByOrderID)
		payments.POST("/webhook", h.ProcessPaymentWebhook)
		payments.POST("/:id/cancel", h.CancelPayment)
		payments.POST("/:id/refund", h.RefundPayment)
	}
}

// CreatePayment creates a new payment
// POST /api/v1/payments
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	// Get user_id from context (set by auth middleware) or header
	userID, _ := c.Get("user_id")
	if userID == nil {
		// Fallback: try from header (for demo/testing)
		if uid := c.GetHeader("X-User-ID"); uid != "" {
			id, _ := strconv.ParseUint(uid, 10, 32)
			userID = uint(id)
		}
	}

	if userID == nil || userID.(uint) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not authenticated",
		})
		return
	}

	var req dto.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	payment, err := h.paymentService.CreatePayment(userID.(uint), &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err == service.ErrPaymentExists {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Payment created successfully",
		"data":    payment,
	})
}

// GetPayment retrieves a payment by ID
// GET /api/v1/payments/:id
func (h *PaymentHandler) GetPayment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid payment ID",
		})
		return
	}

	payment, err := h.paymentService.GetPayment(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    payment,
	})
}

// GetPaymentByOrderID retrieves a payment by order ID
// GET /api/v1/payments/order/:order_id
func (h *PaymentHandler) GetPaymentByOrderID(c *gin.Context) {
	orderID, err := strconv.ParseUint(c.Param("order_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid order ID",
		})
		return
	}

	payment, err := h.paymentService.GetPaymentByOrderID(uint(orderID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    payment,
	})
}

// GetUserPayments retrieves payments for the authenticated user
// GET /api/v1/payments
func (h *PaymentHandler) GetUserPayments(c *gin.Context) {
	userID, _ := c.Get("user_id")
	if userID == nil {
		if uid := c.GetHeader("X-User-ID"); uid != "" {
			id, _ := strconv.ParseUint(uid, 10, 32)
			userID = uint(id)
		}
	}

	if userID == nil || userID.(uint) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not authenticated",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	payments, err := h.paymentService.GetUserPayments(userID.(uint), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    payments,
	})
}

// ProcessPaymentWebhook handles payment provider webhooks
// POST /api/v1/payments/webhook
func (h *PaymentHandler) ProcessPaymentWebhook(c *gin.Context) {
	var req dto.ProcessPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	payment, err := h.paymentService.ProcessPayment(&req)
	if err != nil {
		status := http.StatusInternalServerError
		if err == service.ErrPaymentNotFound {
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
		"message": "Payment processed successfully",
		"data":    payment,
	})
}

// CancelPayment cancels a pending payment
// POST /api/v1/payments/:id/cancel
func (h *PaymentHandler) CancelPayment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid payment ID",
		})
		return
	}

	if err := h.paymentService.CancelPayment(uint(id)); err != nil {
		status := http.StatusInternalServerError
		if err == service.ErrPaymentNotFound {
			status = http.StatusNotFound
		} else if err == service.ErrPaymentNotPending {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Payment cancelled successfully",
	})
}

// RefundPayment processes a refund for a successful payment
// POST /api/v1/payments/:id/refund
func (h *PaymentHandler) RefundPayment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid payment ID",
		})
		return
	}

	var req dto.RefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	if err := h.paymentService.RefundPayment(uint(id), req.Reason); err != nil {
		status := http.StatusInternalServerError
		if err == service.ErrPaymentNotFound {
			status = http.StatusNotFound
		} else if err == service.ErrPaymentNotSuccess {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Payment refunded successfully",
	})
}
