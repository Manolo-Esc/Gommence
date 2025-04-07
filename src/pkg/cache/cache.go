package cache

import (
	"sync"

	"github.com/dgraph-io/ristretto/v2"
)

type CacheService interface {
	Set(key string, value interface{}) bool
	Get(key string) (interface{}, bool)
	Del(key string)
}

type cacheServiceImpl struct {
	provider *ristretto.Cache[string, interface{}]
}

func (c *cacheServiceImpl) Set(key string, value interface{}) bool {
	cost := int64(1)
	c.provider.Set(key, value, cost)
	return true
}

func (c *cacheServiceImpl) Get(key string) (interface{}, bool) {
	value, found := c.provider.Get(key)
	return value, found
}

func (c *cacheServiceImpl) Del(key string) {
	c.provider.Del(key)
}

var (
	cache      *cacheServiceImpl
	createOnce sync.Once
)

// GetCache returns a cache singleton intended to be used across the application
func GetCache() CacheService {
	createOnce.Do(func() {
		var err error
		var provider *ristretto.Cache[string, interface{}]
		provider, err = ristretto.NewCache(&ristretto.Config[string, interface{}]{
			NumCounters: 1e7,     // number of keys to track frequency of (10M).
			MaxCost:     1 << 30, // maximum cost of cache (1GB).
			BufferItems: 64,      // number of keys per Get buffer.
		})
		if err != nil {
			panic(err)
		}
		cache = &cacheServiceImpl{}
		cache.provider = provider
	})
	return cache
}

// NewCache returns a new cache instance. Intended for classes that expect to make a heavy use and don't want to share it
func NewCache() CacheService {
	var err error
	var provider *ristretto.Cache[string, interface{}]
	provider, err = ristretto.NewCache(&ristretto.Config[string, interface{}]{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		panic(err)
	}
	ret := &cacheServiceImpl{}
	ret.provider = provider
	return ret
}

type cacheServiceNopImpl struct {
}

func (c *cacheServiceNopImpl) Set(key string, value interface{}) bool {
	return true
}

func (c *cacheServiceNopImpl) Get(key string) (interface{}, bool) {
	return nil, false
}

func (c *cacheServiceNopImpl) Del(key string) {
}

func GetNopCache() CacheService {
	return &cacheServiceNopImpl{}
}
