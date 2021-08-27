package cache

import (
	"github.com/go-redis/redis/v7"
	"time"
)

type Service interface {
	HSet(key string, values ...interface{}) error
	HSetNX(key string, field string, value interface{}, expiration time.Duration) (set bool, err error)
	Expire(key string, expiration time.Duration) error
	HGet(key string, field string) (string, error)
	HGetAll(key string) map[string]string
	HDel(key string, fields string) error
	Del(key string) error
	Pipeline() redis.Pipeliner
}

func NewCacheService(clusterMode bool, host, port, password string) (Service, error) {
	if clusterMode {
		return NewClusterCacheClient(host, port, password)
	}
	return NewSimpleCacheClient(host, port, password)
}
