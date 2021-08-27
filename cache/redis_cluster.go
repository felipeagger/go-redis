package cache

import (
	"crypto/tls"
	"fmt"
	"github.com/go-redis/redis/v7"
	"strings"
	"sync"
	"time"
)

var (
	onceCacheCluster sync.Once
	cacheClusterClientSvc Service
)

func InitCacheClusterClientSvc(cacheHost string, cachePort string, cachePassword string) {
	onceCacheCluster.Do(func() {

		cacheSvc, err := NewCacheService(true, cacheHost, cachePort, cachePassword)

		if err != nil {
			panic(err)
		}

		cacheClusterClientSvc = cacheSvc
	})
}

func GetCacheClusterClient() Service {
	return cacheClusterClientSvc
}


type ClusterClient struct {
	onceCache sync.Once
	Client    *redis.ClusterClient
}

func NewClusterCacheClient(host, port, password string) (*ClusterClient, error) {

	client, err := NewRedisClusterClient(host, port, password)
	if err != nil {
		return nil, err
	}

	cache := &ClusterClient{
		Client: client,
	}

	return cache, nil
}

func (cache *ClusterClient) HSet(key string, values ...interface{}) error {
	return cache.Client.HSet(key, values).Err()
}

func (cache *ClusterClient) HSetNX(key string, field string, value interface{}, expiration time.Duration) (set bool, err error) {

	set, err = cache.Client.HSetNX(key, field, value).Result()
	if set == true {
		err = cache.Expire(key, expiration)
	}

	return set, err
}

func (cache *ClusterClient) Expire(key string, expiration time.Duration) error {
	return cache.Client.Expire(key, expiration).Err()
}

func (cache *ClusterClient) HGet(key string, field string) (string, error) {

	data, err := cache.Client.HGet(key, field).Result()
	if err != nil && err.Error() == "redis: nil" {
		return data, nil
	}

	return data, err
}

func (cache *ClusterClient) HGetAll(key string) map[string]string {
	return cache.Client.HGetAll(key).Val()
}

func (cache *ClusterClient) HDel(key string, fields string) error {
	return cache.Client.HDel(key, fields).Err()
}

func (cache *ClusterClient) Del(key string) error {
	return cache.Client.Del(key).Err()
}

func (cache *ClusterClient) Pipeline() redis.Pipeliner {
	return cache.Client.Pipeline()
}

func NewRedisClusterClient(redisHost, cachePort, redisPassword string) (cacheClient *redis.ClusterClient, err error) {

	addrs := strings.Split(fmt.Sprintf("%s:%s", redisHost, cachePort), ",")

	if len(redisPassword) > 0 {
		cacheClient = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    addrs,
			Password: redisPassword,
			TLSConfig: &tls.Config{
				RootCAs: nil,
			},
		})
	} else {
		cacheClient = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs: addrs,
		})
	}

	if err := cacheClient.Ping().Err(); err != nil {
		fmt.Println("ERRO NO CLUSTER REDIS")
		return nil, err
	}

	return cacheClient, nil
}
