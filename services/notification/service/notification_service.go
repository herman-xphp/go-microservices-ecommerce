package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
	"time"

	"github.com/herman-xphp/go-microservices-ecommerce/services/notification/domain"
	"github.com/herman-xphp/go-microservices-ecommerce/services/notification/dto"
	"gorm.io/gorm"
)

var (
	ErrTemplateNotFound = errors.New("template not found")
	ErrInvalidRecipient = errors.New("invalid recipient")
)

// EmailSender interface for sending emails
type EmailSender interface {
	Send(to, subject, body string) error
}

// SMSSender interface for sending SMS
type SMSSender interface {
	Send(phoneNumber, message string) error
}

// PushSender interface for sending push notifications
type PushSender interface {
	Send(deviceToken, title, body string, data map[string]string) error
}

// NotificationService defines the interface for notification operations
type NotificationService interface {
	SendEmail(req *dto.SendEmailRequest) (*dto.NotificationResponse, error)
	SendSMS(req *dto.SendSMSRequest) (*dto.NotificationResponse, error)
	SendPush(req *dto.SendPushRequest) (*dto.NotificationResponse, error)
	SendOrderConfirmation(userID uint, email string, data *dto.OrderConfirmationData) error
	SendPaymentSuccess(userID uint, email string, data *dto.PaymentSuccessData) error
	GetUserNotifications(userID uint, page, pageSize int) ([]dto.NotificationResponse, int64, error)
}

type notificationServiceImpl struct {
	db          *gorm.DB
	emailSender EmailSender
	smsSender   SMSSender
	pushSender  PushSender
}

// NewNotificationService creates a new NotificationService
func NewNotificationService(db *gorm.DB, emailSender EmailSender, smsSender SMSSender, pushSender PushSender) NotificationService {
	return &notificationServiceImpl{
		db:          db,
		emailSender: emailSender,
		smsSender:   smsSender,
		pushSender:  pushSender,
	}
}

func (s *notificationServiceImpl) SendEmail(req *dto.SendEmailRequest) (*dto.NotificationResponse, error) {
	body := req.Body

	// Use template if specified
	if req.TemplateID != "" {
		var tmpl domain.NotificationTemplate
		if err := s.db.Where("id = ? AND type = ? AND is_active = true", req.TemplateID, domain.NotificationTypeEmail).First(&tmpl).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, ErrTemplateNotFound
			}
			return nil, err
		}

		rendered, err := s.renderTemplate(tmpl.Body, req.Variables)
		if err != nil {
			return nil, err
		}
		body = rendered
	}

	// Create notification record
	notification := &domain.Notification{
		UserID:     req.UserID,
		Type:       domain.NotificationTypeEmail,
		Status:     domain.NotificationStatusPending,
		Subject:    req.Subject,
		Content:    body,
		Recipient:  req.To,
		TemplateID: req.TemplateID,
	}

	if err := s.db.Create(notification).Error; err != nil {
		return nil, err
	}

	// Send email
	if s.emailSender != nil {
		if err := s.emailSender.Send(req.To, req.Subject, body); err != nil {
			notification.Status = domain.NotificationStatusFailed
			notification.Error = err.Error()
		} else {
			notification.Status = domain.NotificationStatusSent
			now := time.Now()
			notification.SentAt = &now
		}
		s.db.Save(notification)
	} else {
		// Mock: just mark as sent for demo
		notification.Status = domain.NotificationStatusSent
		now := time.Now()
		notification.SentAt = &now
		s.db.Save(notification)
	}

	return s.toNotificationResponse(notification), nil
}

func (s *notificationServiceImpl) SendSMS(req *dto.SendSMSRequest) (*dto.NotificationResponse, error) {
	notification := &domain.Notification{
		UserID:    req.UserID,
		Type:      domain.NotificationTypeSMS,
		Status:    domain.NotificationStatusPending,
		Subject:   "SMS",
		Content:   req.Message,
		Recipient: req.PhoneNumber,
	}

	if err := s.db.Create(notification).Error; err != nil {
		return nil, err
	}

	// Send SMS (mock for now)
	notification.Status = domain.NotificationStatusSent
	now := time.Now()
	notification.SentAt = &now
	s.db.Save(notification)

	return s.toNotificationResponse(notification), nil
}

func (s *notificationServiceImpl) SendPush(req *dto.SendPushRequest) (*dto.NotificationResponse, error) {
	metadata, _ := json.Marshal(req.Data)

	notification := &domain.Notification{
		UserID:    req.UserID,
		Type:      domain.NotificationTypePush,
		Status:    domain.NotificationStatusPending,
		Subject:   req.Title,
		Content:   req.Body,
		Recipient: req.DeviceToken,
		Metadata:  string(metadata),
	}

	if err := s.db.Create(notification).Error; err != nil {
		return nil, err
	}

	// Send push (mock for now)
	notification.Status = domain.NotificationStatusSent
	now := time.Now()
	notification.SentAt = &now
	s.db.Save(notification)

	return s.toNotificationResponse(notification), nil
}

func (s *notificationServiceImpl) SendOrderConfirmation(userID uint, email string, data *dto.OrderConfirmationData) error {
	body := s.buildOrderConfirmationEmail(data)

	_, err := s.SendEmail(&dto.SendEmailRequest{
		UserID:  userID,
		To:      email,
		Subject: "Order Confirmation - #" + string(rune(data.OrderID+'0')),
		Body:    body,
	})
	return err
}

func (s *notificationServiceImpl) SendPaymentSuccess(userID uint, email string, data *dto.PaymentSuccessData) error {
	body := s.buildPaymentSuccessEmail(data)

	_, err := s.SendEmail(&dto.SendEmailRequest{
		UserID:  userID,
		To:      email,
		Subject: "Payment Successful - Transaction " + data.TransactionID,
		Body:    body,
	})
	return err
}

func (s *notificationServiceImpl) GetUserNotifications(userID uint, page, pageSize int) ([]dto.NotificationResponse, int64, error) {
	var notifications []domain.Notification
	var total int64

	s.db.Model(&domain.Notification{}).Where("user_id = ?", userID).Count(&total)

	offset := (page - 1) * pageSize
	if err := s.db.Where("user_id = ?", userID).
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&notifications).Error; err != nil {
		return nil, 0, err
	}

	responses := make([]dto.NotificationResponse, len(notifications))
	for i, n := range notifications {
		responses[i] = *s.toNotificationResponse(&n)
	}

	return responses, total, nil
}

func (s *notificationServiceImpl) renderTemplate(tmpl string, variables map[string]string) (string, error) {
	t, err := template.New("notification").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, variables); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (s *notificationServiceImpl) buildOrderConfirmationEmail(data *dto.OrderConfirmationData) string {
	return `
<!DOCTYPE html>
<html>
<head><style>body{font-family:Arial,sans-serif;}</style></head>
<body>
<h2>Order Confirmation</h2>
<p>Dear ` + data.CustomerName + `,</p>
<p>Thank you for your order! Your order has been confirmed.</p>
<p><strong>Order ID:</strong> #` + string(rune(data.OrderID+'0')) + `</p>
<p><strong>Total Amount:</strong> Rp ` + formatPrice(data.TotalAmount) + `</p>
<p>We will notify you once your order is shipped.</p>
<p>Thank you for shopping with us!</p>
</body>
</html>`
}

func (s *notificationServiceImpl) buildPaymentSuccessEmail(data *dto.PaymentSuccessData) string {
	return `
<!DOCTYPE html>
<html>
<head><style>body{font-family:Arial,sans-serif;}</style></head>
<body>
<h2>Payment Successful</h2>
<p>Your payment has been processed successfully.</p>
<p><strong>Transaction ID:</strong> ` + data.TransactionID + `</p>
<p><strong>Amount:</strong> Rp ` + formatPrice(data.Amount) + `</p>
<p><strong>Payment Method:</strong> ` + data.PaymentMethod + `</p>
<p>Thank you for your purchase!</p>
</body>
</html>`
}

func (s *notificationServiceImpl) toNotificationResponse(n *domain.Notification) *dto.NotificationResponse {
	resp := &dto.NotificationResponse{
		ID:        n.ID,
		UserID:    n.UserID,
		Type:      n.Type,
		Status:    n.Status,
		Subject:   n.Subject,
		Recipient: n.Recipient,
		Error:     n.Error,
		CreatedAt: n.CreatedAt.Format(time.RFC3339),
	}
	if n.SentAt != nil {
		resp.SentAt = n.SentAt.Format(time.RFC3339)
	}
	return resp
}

func formatPrice(price float64) string {
	return string(rune(int(price)))
}
