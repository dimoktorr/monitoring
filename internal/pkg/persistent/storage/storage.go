package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/dimoktorr/monitoring/internal/domain"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"

	"github.com/redis/go-redis/v9"
)

type Storage struct {
	client    redis.UniversalClient
	ttl       time.Duration
	namespace string
	json      jsoniter.API
}

func NewStorage(client redis.UniversalClient, ttl time.Duration, namespace string) *Storage {
	return &Storage{
		client:    client,
		ttl:       ttl,
		namespace: namespace,
		json:      jsoniter.ConfigCompatibleWithStandardLibrary,
	}
}

func (r *Storage) prepareKey(key string) string {
	arr := strings.Split(key, ":")
	if len(arr) == 0 || r.namespace == "" {
		return ""
	}

	if arr[0] != r.namespace {
		return fmt.Sprintf("%s:%s", r.namespace, key)
	}

	return key
}

func (r *Storage) Get(ctx context.Context, key string) (*domain.Product, error) {
	val, err := r.client.Get(ctx, r.prepareKey(key)).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrNoStorageRequest
		}
		return nil, err
	}
	var result *domain.Product

	if err := r.json.Unmarshal([]byte(val), &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *Storage) Set(ctx context.Context, key string, value *domain.Product) error {
	cacheEntries, err := r.json.Marshal(value)
	if err != nil {
		return err
	}

	if err := r.client.Set(ctx, r.prepareKey(key), cacheEntries, r.ttl).Err(); err != nil {
		return fmt.Errorf("save request %s error: %w", key, err)
	}

	return nil
}
