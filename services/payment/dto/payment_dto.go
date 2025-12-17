package dto

import "github.com/herman-xphp/go-microservices-ecommerce/services/payment/domain"

// CreatePaymentRequest represents the payload for creating a payment
type CreatePaymentRequest struct {
	OrderID uint                 `json:"order_id" binding:"required"`
	Amount  float64              `json:"amount" binding:"required,gt=0"`
	Method  domain.PaymentMethod `json:"method" binding:"required"`
}

// PaymentResponse represents a payment in API responses
type PaymentResponse struct {
	ID            uint                 `json:"id"`
	OrderID       uint                 `json:"order_id"`
	UserID        uint                 `json:"user_id"`
	Amount        float64              `json:"amount"`
	Currency      string               `json:"currency"`
	Method        domain.PaymentMethod `json:"method"`
	Status        domain.PaymentStatus `json:"status"`
	TransactionID string               `json:"transaction_id"`
	ProviderRef   string               `json:"provider_ref,omitempty"`
	FailureReason string               `json:"failure_reason,omitempty"`
	PaidAt        string               `json:"paid_at,omitempty"`
	CreatedAt     string               `json:"created_at"`
}

// ProcessPaymentRequest represents a payment processing webhook/callback
type ProcessPaymentRequest struct {
	TransactionID string               `json:"transaction_id" binding:"required"`
	Status        domain.PaymentStatus `json:"status" binding:"required"`
	ProviderRef   string               `json:"provider_ref"`
	FailureReason string               `json:"failure_reason,omitempty"`
}

// PaymentListResponse represents paginated payment list
type PaymentListResponse struct {
	Payments   []PaymentResponse `json:"payments"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}

// RefundRequest represents a refund request
type RefundRequest struct {
	Reason string `json:"reason" binding:"required"`
}
