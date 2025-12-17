package client

import (
	"context"
	"time"

	pb "github.com/herman-xphp/go-microservices-ecommerce/proto/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// AuthClient wraps the gRPC client for Auth Service
type AuthClient struct {
	conn   *grpc.ClientConn
	client pb.AuthServiceClient
}

// UserInfo represents validated user information
type UserInfo struct {
	ID    uint
	Email string
	Name  string
}

// NewAuthClient creates a new gRPC client connection to Auth Service
func NewAuthClient(address string) (*AuthClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	return &AuthClient{
		conn:   conn,
		client: pb.NewAuthServiceClient(conn),
	}, nil
}

// Close closes the gRPC connection
func (c *AuthClient) Close() error {
	return c.conn.Close()
}

// ValidateToken validates a JWT token and returns user info
func (c *AuthClient) ValidateToken(ctx context.Context, token string) (*UserInfo, error) {
	resp, err := c.client.ValidateToken(ctx, &pb.ValidateTokenRequest{
		Token: token,
	})
	if err != nil {
		return nil, err
	}

	if !resp.Valid {
		return nil, nil
	}

	return &UserInfo{
		ID:    uint(resp.UserId),
		Email: resp.Email,
		Name:  "", // Name not available in ValidateToken response
	}, nil
}

// GetUserByID fetches user info by ID
func (c *AuthClient) GetUserByID(ctx context.Context, userID uint) (*UserInfo, error) {
	resp, err := c.client.GetUserById(ctx, &pb.GetUserByIdRequest{
		UserId: uint64(userID),
	})
	if err != nil {
		return nil, err
	}

	if !resp.Found {
		return nil, nil
	}

	return &UserInfo{
		ID:    uint(resp.UserId),
		Email: resp.Email,
		Name:  resp.Name,
	}, nil
}
