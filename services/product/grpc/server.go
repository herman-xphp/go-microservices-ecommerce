package grpc

import (
	"context"
	"errors"

	pb "github.com/herman-xphp/go-microservices-ecommerce/proto/product"
	"github.com/herman-xphp/go-microservices-ecommerce/services/product/service"
)

// ProductGRPCServer implements the gRPC ProductService interface
type ProductGRPCServer struct {
	pb.UnimplementedProductServiceServer
	productService service.ProductService
}

// NewProductGRPCServer creates a new gRPC product server
func NewProductGRPCServer(productService service.ProductService) *ProductGRPCServer {
	return &ProductGRPCServer{productService: productService}
}

// GetProduct returns product info by ID
func (s *ProductGRPCServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	product, err := s.productService.GetProduct(uint(req.ProductId))
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			return &pb.GetProductResponse{Found: false}, nil
		}
		return nil, err
	}

	categoryName := ""
	if product.Category != nil {
		categoryName = product.Category.Name
	}

	return &pb.GetProductResponse{
		Found:        true,
		Id:           uint64(product.ID),
		Name:         product.Name,
		Description:  product.Description,
		Price:        product.Price,
		Stock:        int32(product.Stock),
		CategoryId:   uint64(product.CategoryID),
		CategoryName: categoryName,
		IsActive:     product.IsActive,
	}, nil
}

// CheckStock returns the current stock for a product
func (s *ProductGRPCServer) CheckStock(ctx context.Context, req *pb.CheckStockRequest) (*pb.CheckStockResponse, error) {
	stock, err := s.productService.CheckStock(uint(req.ProductId))
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			return &pb.CheckStockResponse{
				Found:        false,
				ErrorMessage: "product not found",
			}, nil
		}
		return &pb.CheckStockResponse{
			Found:        false,
			ErrorMessage: err.Error(),
		}, nil
	}

	return &pb.CheckStockResponse{
		Found: true,
		Stock: int32(stock),
	}, nil
}

// DecreaseStock reduces the stock for a product
func (s *ProductGRPCServer) DecreaseStock(ctx context.Context, req *pb.DecreaseStockRequest) (*pb.DecreaseStockResponse, error) {
	err := s.productService.DecreaseStock(uint(req.ProductId), int(req.Quantity))
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			return &pb.DecreaseStockResponse{
				Success:      false,
				ErrorMessage: "product not found",
			}, nil
		}
		if errors.Is(err, service.ErrInsufficientStock) {
			return &pb.DecreaseStockResponse{
				Success:      false,
				ErrorMessage: "insufficient stock",
			}, nil
		}
		return &pb.DecreaseStockResponse{
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	// Get remaining stock
	remainingStock, _ := s.productService.CheckStock(uint(req.ProductId))

	return &pb.DecreaseStockResponse{
		Success:        true,
		RemainingStock: int32(remainingStock),
	}, nil
}
