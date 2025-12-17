package client

import (
	"context"
	"time"

	pb "github.com/herman-xphp/go-microservices-ecommerce/proto/product"
	"github.com/herman-xphp/go-microservices-ecommerce/services/cart/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ProductClientImpl implements service.ProductClient using gRPC
type ProductClientImpl struct {
	conn   *grpc.ClientConn
	client pb.ProductServiceClient
}

// NewProductClient creates a new gRPC client for Product Service
func NewProductClient(address string) (*ProductClientImpl, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	return &ProductClientImpl{
		conn:   conn,
		client: pb.NewProductServiceClient(conn),
	}, nil
}

// Close closes the gRPC connection
func (c *ProductClientImpl) Close() error {
	return c.conn.Close()
}

// GetProduct fetches product info by ID
func (c *ProductClientImpl) GetProduct(ctx context.Context, productID uint) (*service.ProductInfo, error) {
	resp, err := c.client.GetProduct(ctx, &pb.GetProductRequest{
		ProductId: uint64(productID),
	})
	if err != nil {
		return nil, err
	}

	if !resp.Found {
		return nil, nil
	}

	return &service.ProductInfo{
		ID:       uint(resp.Id),
		Name:     resp.Name,
		Price:    resp.Price,
		Stock:    int(resp.Stock),
		ImageURL: "", // ImageURL not in proto
		IsActive: resp.IsActive,
	}, nil
}
