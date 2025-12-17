package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/herman-xphp/go-microservices-ecommerce/services/cart/domain"
	"github.com/redis/go-redis/v9"
)

const (
	cartKeyPrefix = "cart:"
	cartTTL       = 7 * 24 * time.Hour // 7 days
)

// CartRepository defines the interface for cart storage
type CartRepository interface {
	Get(ctx context.Context, userID uint) (*domain.Cart, error)
	Save(ctx context.Context, cart *domain.Cart) error
	Delete(ctx context.Context, userID uint) error
}

type redisCartRepository struct {
	client *redis.Client
}

// NewRedisCartRepository creates a new Redis-based cart repository
func NewRedisCartRepository(client *redis.Client) CartRepository {
	return &redisCartRepository{client: client}
}

func (r *redisCartRepository) cartKey(userID uint) string {
	return fmt.Sprintf("%s%d", cartKeyPrefix, userID)
}

func (r *redisCartRepository) Get(ctx context.Context, userID uint) (*domain.Cart, error) {
	key := r.cartKey(userID)

	data, err := r.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		// Return empty cart if not found
		return &domain.Cart{
			UserID: userID,
			Items:  []domain.CartItem{},
		}, nil
	}
	if err != nil {
		return nil, err
	}

	var cart domain.Cart
	if err := json.Unmarshal(data, &cart); err != nil {
		return nil, err
	}

	return &cart, nil
}

func (r *redisCartRepository) Save(ctx context.Context, cart *domain.Cart) error {
	key := r.cartKey(cart.UserID)

	data, err := json.Marshal(cart)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, cartTTL).Err()
}

func (r *redisCartRepository) Delete(ctx context.Context, userID uint) error {
	key := r.cartKey(userID)
	return r.client.Del(ctx, key).Err()
}
