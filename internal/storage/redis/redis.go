package redis

import (
	"context"
	"fmt"
	"github.com/malinatrash/tabhub/internal/config"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	client *redis.Client
}

func New(cfg config.Cache) (*Client, error) {

	addr := fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
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

func (c *Client) Close() {
	err := c.client.Close()
	if err != nil {
		return
	}
}
