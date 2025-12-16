package dao

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"gwh.com/project-project/config"
)

var Rc *RedisCache

type RedisCache struct {
	rdb *redis.Client
}

func init() {
	rdb := redis.NewClient(config.AppConf.InitRedisOptions())
	Rc = &RedisCache{rdb: rdb}
}

func (rc *RedisCache) Put(ctx context.Context, key, value string, ttl time.Duration) error {
	return rc.rdb.Set(ctx, key, value, ttl).Err()
}

func (rc *RedisCache) Get(ctx context.Context, key string) (string, error) {
	result, err := rc.rdb.Get(ctx, key).Result()
	return result, err
}
