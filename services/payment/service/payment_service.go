package service

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/herman-xphp/go-microservices-ecommerce/services/payment/domain"
	"github.com/herman-xphp/go-microservices-ecommerce/services/payment/dto"
	"github.com/herman-xphp/go-microservices-ecommerce/services/payment/repository"
)

var (
	ErrPaymentNotFound   = errors.New("payment not found")
	ErrPaymentExists     = errors.New("payment already exists for this order")
	ErrInvalidStatus     = errors.New("invalid payment status transition")
	ErrPaymentNotPending = errors.New("payment is not in pending status")
	ErrPaymentNotSuccess = errors.New("only successful payments can be refunded")
)

// PaymentService defines the interface for payment operations
type PaymentService interface {
	CreatePayment(userID uint, req *dto.CreatePaymentRequest) (*dto.PaymentResponse, error)
	GetPayment(id uint) (*dto.PaymentResponse, error)
	GetPaymentByOrderID(orderID uint) (*dto.PaymentResponse, error)
	GetUserPayments(userID uint, page, pageSize int) (*dto.PaymentListResponse, error)
	ProcessPayment(req *dto.ProcessPaymentRequest) (*dto.PaymentResponse, error)
	CancelPayment(id uint) error
	RefundPayment(id uint, reason string) error

	// For gRPC
	GetPaymentStatus(orderID uint) (domain.PaymentStatus, error)
	VerifyPayment(orderID uint) (bool, error)
}

type paymentServiceImpl struct {
	paymentRepo repository.PaymentRepository
}

// NewPaymentService creates a new instance of PaymentService
func NewPaymentService(paymentRepo repository.PaymentRepository) PaymentService {
	return &paymentServiceImpl{
		paymentRepo: paymentRepo,
	}
}

func (s *paymentServiceImpl) CreatePayment(userID uint, req *dto.CreatePaymentRequest) (*dto.PaymentResponse, error) {
	// Check if payment already exists for this order
	existing, _ := s.paymentRepo.FindByOrderID(req.OrderID)
	if existing != nil {
		return nil, ErrPaymentExists
	}

	// Generate unique transaction ID
	transactionID := fmt.Sprintf("TXN-%d-%s", time.Now().UnixNano(), uuid.New().String()[:8])

	payment := &domain.Payment{
		OrderID:       req.OrderID,
		UserID:        userID,
		Amount:        req.Amount,
		Currency:      "IDR",
		Method:        req.Method,
		Status:        domain.PaymentStatusPending,
		TransactionID: transactionID,
	}

	if err := s.paymentRepo.Create(payment); err != nil {
		return nil, err
	}

	return s.toPaymentResponse(payment), nil
}

func (s *paymentServiceImpl) GetPayment(id uint) (*dto.PaymentResponse, error) {
	payment, err := s.paymentRepo.FindByID(id)
	if err != nil {
		return nil, ErrPaymentNotFound
	}
	return s.toPaymentResponse(payment), nil
}

func (s *paymentServiceImpl) GetPaymentByOrderID(orderID uint) (*dto.PaymentResponse, error) {
	payment, err := s.paymentRepo.FindByOrderID(orderID)
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return nil, ErrPaymentNotFound
	}
	return s.toPaymentResponse(payment), nil
}

func (s *paymentServiceImpl) GetUserPayments(userID uint, page, pageSize int) (*dto.PaymentListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	payments, total, err := s.paymentRepo.FindByUserID(userID, page, pageSize)
	if err != nil {
		return nil, err
	}

	paymentResponses := make([]dto.PaymentResponse, len(payments))
	for i, payment := range payments {
		paymentResponses[i] = *s.toPaymentResponse(&payment)
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.PaymentListResponse{
		Payments:   paymentResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *paymentServiceImpl) ProcessPayment(req *dto.ProcessPaymentRequest) (*dto.PaymentResponse, error) {
	payment, err := s.paymentRepo.FindByTransactionID(req.TransactionID)
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return nil, ErrPaymentNotFound
	}

	// Update payment based on status
	payment.Status = req.Status
	payment.ProviderRef = req.ProviderRef

	if req.Status == domain.PaymentStatusSuccess {
		now := time.Now()
		payment.PaidAt = &now
	} else if req.Status == domain.PaymentStatusFailed {
		payment.FailureReason = req.FailureReason
	}

	if err := s.paymentRepo.Update(payment); err != nil {
		return nil, err
	}

	return s.toPaymentResponse(payment), nil
}

func (s *paymentServiceImpl) CancelPayment(id uint) error {
	payment, err := s.paymentRepo.FindByID(id)
	if err != nil {
		return ErrPaymentNotFound
	}

	if payment.Status != domain.PaymentStatusPending {
		return ErrPaymentNotPending
	}

	return s.paymentRepo.UpdateStatus(id, domain.PaymentStatusCancelled)
}

func (s *paymentServiceImpl) RefundPayment(id uint, reason string) error {
	payment, err := s.paymentRepo.FindByID(id)
	if err != nil {
		return ErrPaymentNotFound
	}

	if payment.Status != domain.PaymentStatusSuccess {
		return ErrPaymentNotSuccess
	}

	payment.Status = domain.PaymentStatusRefunded
	payment.FailureReason = "Refund: " + reason

	return s.paymentRepo.Update(payment)
}

func (s *paymentServiceImpl) GetPaymentStatus(orderID uint) (domain.PaymentStatus, error) {
	payment, err := s.paymentRepo.FindByOrderID(orderID)
	if err != nil {
		return "", err
	}
	if payment == nil {
		return "", ErrPaymentNotFound
	}
	return payment.Status, nil
}

func (s *paymentServiceImpl) VerifyPayment(orderID uint) (bool, error) {
	payment, err := s.paymentRepo.FindByOrderID(orderID)
	if err != nil {
		return false, err
	}
	if payment == nil {
		return false, nil
	}
	return payment.Status == domain.PaymentStatusSuccess, nil
}

// Helper: convert domain.Payment to dto.PaymentResponse
func (s *paymentServiceImpl) toPaymentResponse(payment *domain.Payment) *dto.PaymentResponse {
	resp := &dto.PaymentResponse{
		ID:            payment.ID,
		OrderID:       payment.OrderID,
		UserID:        payment.UserID,
		Amount:        payment.Amount,
		Currency:      payment.Currency,
		Method:        payment.Method,
		Status:        payment.Status,
		TransactionID: payment.TransactionID,
		ProviderRef:   payment.ProviderRef,
		FailureReason: payment.FailureReason,
		CreatedAt:     payment.CreatedAt.Format(time.RFC3339),
	}

	if payment.PaidAt != nil {
		resp.PaidAt = payment.PaidAt.Format(time.RFC3339)
	}

	return resp
}
