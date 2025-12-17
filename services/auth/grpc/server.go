package grpc

import (
	"context"

	pb "github.com/username/go-microservices-ecommerce/proto/auth"
	"github.com/username/go-microservices-ecommerce/services/auth/service"
)

// AuthGRPCServer implements the gRPC AuthService interface
type AuthGRPCServer struct {
	pb.UnimplementedAuthServiceServer
	authService service.AuthService
}

// NewAuthGRPCServer creates a new gRPC auth server
func NewAuthGRPCServer(authService service.AuthService) *AuthGRPCServer {
	return &AuthGRPCServer{authService: authService}
}

// ValidateToken validates a JWT token and returns user info
func (s *AuthGRPCServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	user, err := s.authService.ValidateToken(req.Token)
	if err != nil {
		return &pb.ValidateTokenResponse{
			Valid:        false,
			ErrorMessage: err.Error(),
		}, nil
	}

	return &pb.ValidateTokenResponse{
		Valid:  true,
		UserId: uint64(user.ID),
		Email:  user.Email,
		Role:   user.Role,
	}, nil
}

// GetUserById returns user information by ID
func (s *AuthGRPCServer) GetUserById(ctx context.Context, req *pb.GetUserByIdRequest) (*pb.GetUserByIdResponse, error) {
	user, err := s.authService.GetUserByID(uint(req.UserId))
	if err != nil {
		return &pb.GetUserByIdResponse{
			Found: false,
		}, nil
	}

	return &pb.GetUserByIdResponse{
		Found:  true,
		UserId: uint64(user.ID),
		Email:  user.Email,
		Name:   user.Name,
		Role:   user.Role,
	}, nil
}
