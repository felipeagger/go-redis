package cache

import (
	"crypto/tls"
	"fmt"
	"github.com/go-redis/redis/v7"
	"sync"
	"time"
)

var (
	onceCache sync.Once
	cacheClientSvc Service
	)

func InitCacheClientSvc(cacheHost string, cachePort string, cachePassword string) {
	onceCache.Do(func() {

		cacheSvc, err := NewCacheService(false, cacheHost, cachePort, cachePassword)

		if err != nil {
			panic(err)
		}

		cacheClientSvc = cacheSvc
	})
}

func GetCacheClient() Service {
	return cacheClientSvc
}


type Client struct {
	Client    *redis.Client
}

func NewSimpleCacheClient(host, port, password string) (*Client, error) {

	client, err := NewRedisClient(host, port, password)
	if err != nil {
		return nil, err
	}

	cache := &Client{
		Client: client,
	}

	return cache, nil
}

func (cache *Client) HSet(key string, values ...interface{}) error {
	return cache.Client.HSet(key, values).Err()
}

func (cache *Client) HSetNX(key string, field string, value interface{}, expiration time.Duration) (set bool, err error) {

	set, err = cache.Client.HSetNX(key, field, value).Result()
	if set == true {
		err = cache.Expire(key, expiration)
	}

	return set, err
}

func (cache *Client) Expire(key string, expiration time.Duration) error {
	return cache.Client.Expire(key, expiration).Err()
}

func (cache *Client) HGet(key string, field string) (string, error) {

	data, err := cache.Client.HGet(key, field).Result()
	if err != nil && err.Error() == "redis: nil" {
		return data, nil
	}

	return data, err
}

func (cache *Client) HGetAll(key string) map[string]string {
	return cache.Client.HGetAll(key).Val()
}

func (cache *Client) HDel(key string, fields string) error {
	return cache.Client.HDel(key, fields).Err()
}

func (cache *Client) Del(key string) error {
	return cache.Client.Del(key).Err()
}

func (cache *Client) Pipeline() redis.Pipeliner {
	return cache.Client.Pipeline()
}

// NewCacheClient return a new instance of cache client
func NewRedisClient(hostname, port, password string) (*redis.Client, error) {

	var client *redis.Client

	cachePort := "6379"
	if port != "" {
		cachePort = port
	}

	if len(password) > 0 {
		client = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", hostname, cachePort),
			DB:       0, // use default DB
			Password: password,
			TLSConfig: &tls.Config{
				RootCAs: nil,
			},
		})
	} else {
		client = redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%s", hostname, cachePort),
			DB:   0, // use default DB
		})
	}

	_, err := client.Ping().Result()

	if err != nil {
		println("ERROR ON REDIS: NewCacheClient()")
		return nil, err
	}

	return client, nil
}
