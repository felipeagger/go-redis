package main

import (
	"github.com/felipeagger/go-redis/cache"
)


func init() {
	cache.InitCacheClientSvc("0.0.0.0", "6379", "")

	cache.InitCacheClusterClientSvc("0.0.0.0", "7005", "")
}

func main() {
  
  cache.GetCacheClient().HSet("test", "key", "value")
  cache.GetCacheClusterClient().HSet("test", "key", "value")
  
  cache.GetCacheClient().HGet("test", "key")
  cache.GetCacheClusterClient().HGet("test", "key")

}