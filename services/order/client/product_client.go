package client

import (
	"context"
	"log"
	"time"

	pb "github.com/herman-xphp/go-microservices-ecommerce/proto/product"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ProductClient wraps the gRPC client for Product Service
type ProductClient struct {
	conn   *grpc.ClientConn
	client pb.ProductServiceClient
}

// ProductInfo represents product data returned from Product Service
type ProductInfo struct {
	ID          uint
	Name        string
	Description string
	Price       float64
	Stock       int
	IsActive    bool
}

// NewProductClient creates a new gRPC client connection to Product Service
func NewProductClient(address string) (*ProductClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	log.Printf("✅ Connected to Product Service at %s", address)
	return &ProductClient{
		conn:   conn,
		client: pb.NewProductServiceClient(conn),
	}, nil
}

// Close closes the gRPC connection
func (c *ProductClient) Close() error {
	return c.conn.Close()
}

// GetProduct fetches product info by ID
func (c *ProductClient) GetProduct(ctx context.Context, productID uint) (*ProductInfo, error) {
	resp, err := c.client.GetProduct(ctx, &pb.GetProductRequest{
		ProductId: uint64(productID),
	})
	if err != nil {
		return nil, err
	}

	if !resp.Found {
		return nil, nil
	}

	return &ProductInfo{
		ID:          uint(resp.Id),
		Name:        resp.Name,
		Description: resp.Description,
		Price:       resp.Price,
		Stock:       int(resp.Stock),
		IsActive:    resp.IsActive,
	}, nil
}

// CheckStock checks stock availability for a product
func (c *ProductClient) CheckStock(ctx context.Context, productID uint) (int, error) {
	resp, err := c.client.CheckStock(ctx, &pb.CheckStockRequest{
		ProductId: uint64(productID),
	})
	if err != nil {
		return 0, err
	}

	if !resp.Found {
		return 0, nil
	}

	return int(resp.Stock), nil
}

// DecreaseStock decreases stock for a product (called when order is confirmed)
func (c *ProductClient) DecreaseStock(ctx context.Context, productID uint, quantity int) (int, error) {
	resp, err := c.client.DecreaseStock(ctx, &pb.DecreaseStockRequest{
		ProductId: uint64(productID),
		Quantity:  int32(quantity),
	})
	if err != nil {
		return 0, err
	}

	if !resp.Success {
		return 0, &StockError{Message: resp.ErrorMessage}
	}

	return int(resp.RemainingStock), nil
}

// StockError represents a stock-related error
type StockError struct {
	Message string
}

func (e *StockError) Error() string {
	return e.Message
}
