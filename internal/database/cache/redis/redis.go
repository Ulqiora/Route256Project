package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Ulqiora/Route256Project/internal/config"
	"github.com/redis/go-redis/v9"
)

var defaultTimeout = 10 * time.Minute

type Redis struct {
	redis *redis.Client
}

func New(configRedis config.RedisConfig) *Redis {
	return &Redis{
		redis: redis.NewClient(&redis.Options{
			Addr:     configRedis.Address,
			Password: configRedis.Password,
			DB:       configRedis.DB,
		}),
	}
}

func (r *Redis) Ping(ctx context.Context) error {
	return r.redis.Ping(ctx).Err()
}

func (r *Redis) Get(ctx context.Context, key string) ([]byte, error) {
	result, err := r.redis.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, errors.New("key not found")
		}
		return nil, fmt.Errorf("failed get cached data from redis : %w", err)
	}
	return []byte(result), nil
}

func (r *Redis) Set(ctx context.Context, key string, a any) error {
	data, err := json.Marshal(a)
	if err != nil {
		return err
	}
	errRedis := r.redis.Set(ctx, key, data, defaultTimeout)
	return errRedis.Err()
}

func (r *Redis) SetWithTTL(ctx context.Context, timeout time.Duration, key string, a any) error {
	data, err := json.Marshal(a)
	if err != nil {
		return err
	}
	errRedis := r.redis.Set(ctx, key, data, timeout)
	return errRedis.Err()
}

func (r *Redis) SetMulti(ctx context.Context, a map[string]any) error {
	for key, object := range a {
		data, err := json.Marshal(object)
		if err != nil {
			return err
		}
		errRedis := r.redis.Set(ctx, key, data, defaultTimeout)
		return errRedis.Err()
	}
	return nil
}

func (r *Redis) SetMultiWithTTL(ctx context.Context, timeout time.Duration, a map[string]any) error {
	for key, object := range a {
		data, err := json.Marshal(object)
		if err != nil {
			return err
		}
		errRedis := r.redis.Set(ctx, key, data, timeout)
		return errRedis.Err()
	}
	return nil
}
