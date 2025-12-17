package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/herman-xphp/go-microservices-ecommerce/services/notification/dto"
	"github.com/herman-xphp/go-microservices-ecommerce/services/notification/service"
)

// NotificationHandler handles HTTP requests for notifications
type NotificationHandler struct {
	notificationService service.NotificationService
}

// NewNotificationHandler creates a new NotificationHandler
func NewNotificationHandler(notificationService service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: notificationService}
}

// RegisterRoutes registers notification routes
func (h *NotificationHandler) RegisterRoutes(router *gin.RouterGroup) {
	notifications := router.Group("/notifications")
	{
		notifications.POST("/email", h.SendEmail)
		notifications.POST("/sms", h.SendSMS)
		notifications.POST("/push", h.SendPush)
		notifications.GET("", h.GetUserNotifications)
	}
}

// SendEmail sends an email notification
// POST /api/v1/notifications/email
func (h *NotificationHandler) SendEmail(c *gin.Context) {
	var req dto.SendEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	notification, err := h.notificationService.SendEmail(&req)
	if err != nil {
		status := http.StatusInternalServerError
		if err == service.ErrTemplateNotFound {
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
		"message": "Email sent successfully",
		"data":    notification,
	})
}

// SendSMS sends an SMS notification
// POST /api/v1/notifications/sms
func (h *NotificationHandler) SendSMS(c *gin.Context) {
	var req dto.SendSMSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	notification, err := h.notificationService.SendSMS(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "SMS sent successfully",
		"data":    notification,
	})
}

// SendPush sends a push notification
// POST /api/v1/notifications/push
func (h *NotificationHandler) SendPush(c *gin.Context) {
	var req dto.SendPushRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	notification, err := h.notificationService.SendPush(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Push notification sent successfully",
		"data":    notification,
	})
}

// GetUserNotifications retrieves notifications for a user
// GET /api/v1/notifications
func (h *NotificationHandler) GetUserNotifications(c *gin.Context) {
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

	notifications, total, err := h.notificationService.GetUserNotifications(userID.(uint), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"notifications": notifications,
			"total":         total,
			"page":          page,
			"page_size":     pageSize,
		},
	})
}
