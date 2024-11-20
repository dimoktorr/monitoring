package storage

import (
	"context"
	"strings"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

func NewRedisUniversalClient(ctx context.Context, cfg Config) (redis.UniversalClient, error) {
	addr := strings.Split(cfg.Hosts, ",")
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    addr,
		Password: cfg.Password,
	})

	// Enable tracing instrumentation.
	if err := redisotel.InstrumentTracing(client); err != nil {
		return nil, err
	}

	return client, nil
}
