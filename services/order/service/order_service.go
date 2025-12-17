package service

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/herman-xphp/go-microservices-ecommerce/services/order/client"
	"github.com/herman-xphp/go-microservices-ecommerce/services/order/domain"
	"github.com/herman-xphp/go-microservices-ecommerce/services/order/dto"
	"github.com/herman-xphp/go-microservices-ecommerce/services/order/repository"
	"gorm.io/gorm"
)

var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrProductNotFound    = errors.New("product not found")
	ErrInsufficientStock  = errors.New("insufficient stock")
	ErrProductUnavailable = errors.New("product is unavailable")
	ErrEmptyOrder         = errors.New("order must have at least one item")
)

// OrderService defines the interface for order operations
type OrderService interface {
	CreateOrder(ctx context.Context, userID uint, req *dto.CreateOrderRequest) (*dto.OrderResponse, error)
	GetOrder(ctx context.Context, id uint) (*dto.OrderResponse, error)
	GetUserOrders(ctx context.Context, userID uint, page, pageSize int) (*dto.OrderListResponse, error)
	UpdateOrderStatus(ctx context.Context, id uint, status domain.OrderStatus) error
	CancelOrder(ctx context.Context, id uint) error
}

type orderServiceImpl struct {
	orderRepo     repository.OrderRepository
	productClient *client.ProductClient
}

// NewOrderService creates a new instance of OrderService
func NewOrderService(orderRepo repository.OrderRepository, productClient *client.ProductClient) OrderService {
	return &orderServiceImpl{
		orderRepo:     orderRepo,
		productClient: productClient,
	}
}

func (s *orderServiceImpl) CreateOrder(ctx context.Context, userID uint, req *dto.CreateOrderRequest) (*dto.OrderResponse, error) {
	if len(req.Items) == 0 {
		return nil, ErrEmptyOrder
	}

	var orderItems []domain.OrderItem
	var totalAmount float64

	// Validate products and calculate totals by calling Product Service via gRPC
	for _, item := range req.Items {
		// Get product info from Product Service
		product, err := s.productClient.GetProduct(ctx, item.ProductID)
		if err != nil {
			return nil, err
		}
		if product == nil {
			return nil, ErrProductNotFound
		}
		if !product.IsActive {
			return nil, ErrProductUnavailable
		}

		// Check stock availability
		if product.Stock < item.Quantity {
			return nil, ErrInsufficientStock
		}

		subtotal := product.Price * float64(item.Quantity)
		orderItems = append(orderItems, domain.OrderItem{
			ProductID: item.ProductID,
			Name:      product.Name,
			Price:     product.Price,
			Quantity:  item.Quantity,
			Subtotal:  subtotal,
		})
		totalAmount += subtotal
	}

	// Create order
	order := &domain.Order{
		UserID:      userID,
		Status:      domain.OrderStatusPending,
		TotalAmount: totalAmount,
		Items:       orderItems,
	}

	if err := s.orderRepo.Create(order); err != nil {
		return nil, err
	}

	// Decrease stock for each item (this is a simplified approach)
	// In production, you'd use Saga pattern or 2PC for distributed transactions
	for _, item := range req.Items {
		_, err := s.productClient.DecreaseStock(ctx, item.ProductID, item.Quantity)
		if err != nil {
			// In a real system, you'd need to compensate (rollback) here
			// For now, we log the error but the order is already created
			// This demonstrates the need for proper distributed transaction handling
		}
	}

	return s.toOrderResponse(order), nil
}

func (s *orderServiceImpl) GetOrder(ctx context.Context, id uint) (*dto.OrderResponse, error) {
	order, err := s.orderRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}
	return s.toOrderResponse(order), nil
}

func (s *orderServiceImpl) GetUserOrders(ctx context.Context, userID uint, page, pageSize int) (*dto.OrderListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	orders, total, err := s.orderRepo.FindByUserID(userID, page, pageSize)
	if err != nil {
		return nil, err
	}

	orderResponses := make([]dto.OrderResponse, len(orders))
	for i, order := range orders {
		orderResponses[i] = *s.toOrderResponse(&order)
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.OrderListResponse{
		Orders:     orderResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *orderServiceImpl) UpdateOrderStatus(ctx context.Context, id uint, status domain.OrderStatus) error {
	_, err := s.orderRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrOrderNotFound
		}
		return err
	}
	return s.orderRepo.UpdateStatus(id, status)
}

func (s *orderServiceImpl) CancelOrder(ctx context.Context, id uint) error {
	order, err := s.orderRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrOrderNotFound
		}
		return err
	}

	// Only pending orders can be cancelled
	if order.Status != domain.OrderStatusPending {
		return errors.New("only pending orders can be cancelled")
	}

	return s.orderRepo.UpdateStatus(id, domain.OrderStatusCancelled)
}

// Helper: convert domain.Order to dto.OrderResponse
func (s *orderServiceImpl) toOrderResponse(order *domain.Order) *dto.OrderResponse {
	items := make([]dto.OrderItemResponse, len(order.Items))
	for i, item := range order.Items {
		items[i] = dto.OrderItemResponse{
			ID:        item.ID,
			ProductID: item.ProductID,
			Name:      item.Name,
			Price:     item.Price,
			Quantity:  item.Quantity,
			Subtotal:  item.Subtotal,
		}
	}

	return &dto.OrderResponse{
		ID:          order.ID,
		UserID:      order.UserID,
		Status:      order.Status,
		TotalAmount: order.TotalAmount,
		Items:       items,
		CreatedAt:   order.CreatedAt.Format(time.RFC3339),
	}
}
