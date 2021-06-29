package storage

import (
	"context"

	"github.com/gxravel/bus-routes/internal/config"

	"github.com/go-redis/redis/v8"
)

// Client is a client for interaction with storage.
type Client struct {
	*redis.Client
}

// NewClient creates new instance of Client.
func NewClient(cfg config.Storage) (*Client, error) {
	cli := redis.NewClient(&redis.Options{
		Addr: cfg.RedisDSN,
	})
	ctx := context.Background()
	_, err := cli.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &Client{
		cli,
	}, nil
}
