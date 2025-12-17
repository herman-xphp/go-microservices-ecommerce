package domain

import "time"

// NotificationType represents the type of notification
type NotificationType string

const (
	NotificationTypeEmail NotificationType = "email"
	NotificationTypeSMS   NotificationType = "sms"
	NotificationTypePush  NotificationType = "push"
)

// NotificationStatus represents the delivery status
type NotificationStatus string

const (
	NotificationStatusPending   NotificationStatus = "pending"
	NotificationStatusSent      NotificationStatus = "sent"
	NotificationStatusFailed    NotificationStatus = "failed"
	NotificationStatusDelivered NotificationStatus = "delivered"
)

// Notification represents a notification record
type Notification struct {
	ID         uint               `json:"id" gorm:"primaryKey"`
	UserID     uint               `json:"user_id" gorm:"index"`
	Type       NotificationType   `json:"type"`
	Status     NotificationStatus `json:"status" gorm:"default:pending"`
	Subject    string             `json:"subject"`
	Content    string             `json:"content" gorm:"type:text"`
	Recipient  string             `json:"recipient"` // email/phone/device_token
	TemplateID string             `json:"template_id,omitempty"`
	Metadata   string             `json:"metadata,omitempty" gorm:"type:text"` // JSON metadata
	SentAt     *time.Time         `json:"sent_at"`
	Error      string             `json:"error,omitempty"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
}

// TableName overrides the table name
func (Notification) TableName() string {
	return "notifications"
}

// NotificationTemplate represents email/sms templates
type NotificationTemplate struct {
	ID        string           `json:"id" gorm:"primaryKey"` // e.g., "order_confirmation"
	Type      NotificationType `json:"type"`
	Subject   string           `json:"subject"`
	Body      string           `json:"body" gorm:"type:text"`
	IsActive  bool             `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

func (NotificationTemplate) TableName() string {
	return "notification_templates"
}
