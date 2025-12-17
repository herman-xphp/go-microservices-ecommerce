package domain

import "time"

// Order represents an order in the system
type Order struct {
	ID          uint        `json:"id" gorm:"primaryKey"`
	UserID      uint        `json:"user_id" gorm:"not null;index"`
	Status      OrderStatus `json:"status" gorm:"default:pending"`
	TotalAmount float64     `json:"total_amount" gorm:"not null"`
	Items       []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// TableName overrides the table name
func (Order) TableName() string {
	return "orders"
}

// OrderItem represents a single item in an order
type OrderItem struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	OrderID   uint    `json:"order_id" gorm:"not null;index"`
	ProductID uint    `json:"product_id" gorm:"not null"`
	Name      string  `json:"name" gorm:"not null"`
	Price     float64 `json:"price" gorm:"not null"`
	Quantity  int     `json:"quantity" gorm:"not null"`
	Subtotal  float64 `json:"subtotal" gorm:"not null"`
}

// TableName overrides the table name
func (OrderItem) TableName() string {
	return "order_items"
}

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
)
