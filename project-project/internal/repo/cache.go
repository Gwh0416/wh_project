package repo

import (
	"context"
	"time"
)

type Cache interface {
	Put(ctx context.Context, key, value string, time time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	HKeys(ctx context.Context, key string) ([]string, error)
	Delete(background context.Context, keys []string)
	HSet(ctx context.Context, key string, field string, value string)
}
