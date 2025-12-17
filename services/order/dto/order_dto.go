package dto

import "github.com/herman-xphp/go-microservices-ecommerce/services/order/domain"

// CreateOrderRequest represents the payload for creating an order
type CreateOrderRequest struct {
	Items []OrderItemRequest `json:"items" binding:"required,min=1,dive"`
}

// OrderItemRequest represents a single item in the order request
type OrderItemRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

// OrderResponse represents an order in API responses
type OrderResponse struct {
	ID          uint                `json:"id"`
	UserID      uint                `json:"user_id"`
	Status      domain.OrderStatus  `json:"status"`
	TotalAmount float64             `json:"total_amount"`
	Items       []OrderItemResponse `json:"items"`
	CreatedAt   string              `json:"created_at"`
}

// OrderItemResponse represents an order item in API responses
type OrderItemResponse struct {
	ID        uint    `json:"id"`
	ProductID uint    `json:"product_id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
	Subtotal  float64 `json:"subtotal"`
}

// UpdateOrderStatusRequest represents the payload for updating order status
type UpdateOrderStatusRequest struct {
	Status domain.OrderStatus `json:"status" binding:"required"`
}

// OrderListResponse represents paginated order list
type OrderListResponse struct {
	Orders     []OrderResponse `json:"orders"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
}
