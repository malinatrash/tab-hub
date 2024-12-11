package redis

import (
	"context"
	"github.com/malinatrash/tabhub/internal/config"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	client *redis.Client
}

func New(cfg config.Cache) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Username: cfg.Username,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Client{client: rdb}, nil
}
