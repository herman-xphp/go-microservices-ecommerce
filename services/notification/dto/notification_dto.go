package dto

import "github.com/herman-xphp/go-microservices-ecommerce/services/notification/domain"

// SendEmailRequest represents a request to send an email
type SendEmailRequest struct {
	UserID     uint              `json:"user_id" binding:"required"`
	To         string            `json:"to" binding:"required,email"`
	Subject    string            `json:"subject" binding:"required"`
	Body       string            `json:"body"`
	TemplateID string            `json:"template_id,omitempty"`
	Variables  map[string]string `json:"variables,omitempty"`
}

// SendSMSRequest represents a request to send an SMS
type SendSMSRequest struct {
	UserID      uint              `json:"user_id" binding:"required"`
	PhoneNumber string            `json:"phone_number" binding:"required"`
	Message     string            `json:"message" binding:"required"`
	TemplateID  string            `json:"template_id,omitempty"`
	Variables   map[string]string `json:"variables,omitempty"`
}

// SendPushRequest represents a request to send a push notification
type SendPushRequest struct {
	UserID      uint              `json:"user_id" binding:"required"`
	DeviceToken string            `json:"device_token" binding:"required"`
	Title       string            `json:"title" binding:"required"`
	Body        string            `json:"body" binding:"required"`
	Data        map[string]string `json:"data,omitempty"`
}

// NotificationResponse represents a notification in responses
type NotificationResponse struct {
	ID        uint                      `json:"id"`
	UserID    uint                      `json:"user_id"`
	Type      domain.NotificationType   `json:"type"`
	Status    domain.NotificationStatus `json:"status"`
	Subject   string                    `json:"subject"`
	Recipient string                    `json:"recipient"`
	SentAt    string                    `json:"sent_at,omitempty"`
	Error     string                    `json:"error,omitempty"`
	CreatedAt string                    `json:"created_at"`
}

// OrderConfirmationData represents data for order confirmation email
type OrderConfirmationData struct {
	OrderID      uint            `json:"order_id"`
	CustomerName string          `json:"customer_name"`
	TotalAmount  float64         `json:"total_amount"`
	OrderItems   []OrderItemData `json:"order_items"`
}

type OrderItemData struct {
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}

// PaymentSuccessData represents data for payment success email
type PaymentSuccessData struct {
	OrderID       uint    `json:"order_id"`
	PaymentID     uint    `json:"payment_id"`
	Amount        float64 `json:"amount"`
	PaymentMethod string  `json:"payment_method"`
	TransactionID string  `json:"transaction_id"`
}
