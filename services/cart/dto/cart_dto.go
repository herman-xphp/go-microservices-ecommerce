package dto

// AddToCartRequest represents adding an item to cart
type AddToCartRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

// UpdateCartItemRequest represents updating cart item quantity
type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=0"`
}

// CartItemResponse represents a cart item in responses
type CartItemResponse struct {
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	Subtotal    float64 `json:"subtotal"`
	ImageURL    string  `json:"image_url,omitempty"`
}

// CartResponse represents the cart in API responses
type CartResponse struct {
	UserID     uint               `json:"user_id"`
	Items      []CartItemResponse `json:"items"`
	TotalItems int                `json:"total_items"`
	TotalPrice float64            `json:"total_price"`
}

// CheckoutRequest represents a checkout request
type CheckoutRequest struct {
	ShippingAddress string `json:"shipping_address" binding:"required"`
	PaymentMethod   string `json:"payment_method" binding:"required"`
}
