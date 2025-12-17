package service

import (
	"context"
	"errors"

	"github.com/herman-xphp/go-microservices-ecommerce/services/cart/domain"
	"github.com/herman-xphp/go-microservices-ecommerce/services/cart/dto"
	"github.com/herman-xphp/go-microservices-ecommerce/services/cart/repository"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrCartEmpty       = errors.New("cart is empty")
	ErrItemNotInCart   = errors.New("item not in cart")
)

// ProductInfo represents product info from Product Service
type ProductInfo struct {
	ID       uint
	Name     string
	Price    float64
	Stock    int
	ImageURL string
	IsActive bool
}

// ProductClient interface for getting product info
type ProductClient interface {
	GetProduct(ctx context.Context, productID uint) (*ProductInfo, error)
}

// CartService defines the interface for cart operations
type CartService interface {
	GetCart(ctx context.Context, userID uint) (*dto.CartResponse, error)
	AddToCart(ctx context.Context, userID uint, req *dto.AddToCartRequest) (*dto.CartResponse, error)
	UpdateItem(ctx context.Context, userID uint, productID uint, req *dto.UpdateCartItemRequest) (*dto.CartResponse, error)
	RemoveItem(ctx context.Context, userID uint, productID uint) (*dto.CartResponse, error)
	ClearCart(ctx context.Context, userID uint) error
}

type cartServiceImpl struct {
	cartRepo      repository.CartRepository
	productClient ProductClient
}

// NewCartService creates a new CartService
func NewCartService(cartRepo repository.CartRepository, productClient ProductClient) CartService {
	return &cartServiceImpl{
		cartRepo:      cartRepo,
		productClient: productClient,
	}
}

func (s *cartServiceImpl) GetCart(ctx context.Context, userID uint) (*dto.CartResponse, error) {
	cart, err := s.cartRepo.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.toCartResponse(cart), nil
}

func (s *cartServiceImpl) AddToCart(ctx context.Context, userID uint, req *dto.AddToCartRequest) (*dto.CartResponse, error) {
	// Get product info from Product Service
	product, err := s.productClient.GetProduct(ctx, req.ProductID)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, ErrProductNotFound
	}

	// Get current cart
	cart, err := s.cartRepo.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Add item to cart
	cart.AddItem(domain.CartItem{
		ProductID:   product.ID,
		ProductName: product.Name,
		Price:       product.Price,
		Quantity:    req.Quantity,
		ImageURL:    product.ImageURL,
	})

	// Save cart
	if err := s.cartRepo.Save(ctx, cart); err != nil {
		return nil, err
	}

	return s.toCartResponse(cart), nil
}

func (s *cartServiceImpl) UpdateItem(ctx context.Context, userID uint, productID uint, req *dto.UpdateCartItemRequest) (*dto.CartResponse, error) {
	cart, err := s.cartRepo.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	if !cart.UpdateItemQuantity(productID, req.Quantity) {
		return nil, ErrItemNotInCart
	}

	if err := s.cartRepo.Save(ctx, cart); err != nil {
		return nil, err
	}

	return s.toCartResponse(cart), nil
}

func (s *cartServiceImpl) RemoveItem(ctx context.Context, userID uint, productID uint) (*dto.CartResponse, error) {
	cart, err := s.cartRepo.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	if !cart.RemoveItem(productID) {
		return nil, ErrItemNotInCart
	}

	if err := s.cartRepo.Save(ctx, cart); err != nil {
		return nil, err
	}

	return s.toCartResponse(cart), nil
}

func (s *cartServiceImpl) ClearCart(ctx context.Context, userID uint) error {
	return s.cartRepo.Delete(ctx, userID)
}

func (s *cartServiceImpl) toCartResponse(cart *domain.Cart) *dto.CartResponse {
	items := make([]dto.CartItemResponse, len(cart.Items))
	for i, item := range cart.Items {
		items[i] = dto.CartItemResponse{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Price:       item.Price,
			Quantity:    item.Quantity,
			Subtotal:    item.Price * float64(item.Quantity),
			ImageURL:    item.ImageURL,
		}
	}

	return &dto.CartResponse{
		UserID:     cart.UserID,
		Items:      items,
		TotalItems: cart.TotalItems,
		TotalPrice: cart.TotalPrice,
	}
}
