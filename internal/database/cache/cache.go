package cache

import (
	"context"
	"time"
)

type Cache interface {
	Ping(ctx context.Context) error
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, a any) error
	SetWithTTL(ctx context.Context, timeout time.Duration, key string, a any) error
	SetMulti(ctx context.Context, a map[string]any) error
	SetMultiWithTTL(ctx context.Context, timeout time.Duration, a map[string]any) error
}
