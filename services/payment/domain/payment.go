package domain

import "time"

// Payment represents a payment transaction
type Payment struct {
	ID            uint          `json:"id" gorm:"primaryKey"`
	OrderID       uint          `json:"order_id" gorm:"not null;index"`
	UserID        uint          `json:"user_id" gorm:"not null;index"`
	Amount        float64       `json:"amount" gorm:"not null"`
	Currency      string        `json:"currency" gorm:"default:IDR"`
	Method        PaymentMethod `json:"method" gorm:"not null"`
	Status        PaymentStatus `json:"status" gorm:"default:pending"`
	TransactionID string        `json:"transaction_id" gorm:"uniqueIndex"`
	ProviderRef   string        `json:"provider_ref"` // Reference from payment provider
	FailureReason string        `json:"failure_reason"`
	PaidAt        *time.Time    `json:"paid_at"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

// TableName overrides the table name
func (Payment) TableName() string {
	return "payments"
}

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "pending"
	PaymentStatusProcessing PaymentStatus = "processing"
	PaymentStatusSuccess    PaymentStatus = "success"
	PaymentStatusFailed     PaymentStatus = "failed"
	PaymentStatusRefunded   PaymentStatus = "refunded"
	PaymentStatusCancelled  PaymentStatus = "cancelled"
)

// PaymentMethod represents supported payment methods
type PaymentMethod string

const (
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
	PaymentMethodCreditCard   PaymentMethod = "credit_card"
	PaymentMethodEWallet      PaymentMethod = "e_wallet"
	PaymentMethodVA           PaymentMethod = "virtual_account"
	PaymentMethodQRIS         PaymentMethod = "qris"
)
