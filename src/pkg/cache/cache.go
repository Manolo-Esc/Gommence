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

/*
func f() {
	cache, err := ristretto.NewCache(&ristretto.Config[string, interface{}]{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		panic(err)
	}
	defer cache.Close()

	// set a value with a cost of 1
	cache.Set("key", "value", 1)
	cache.Set("key", 5, 1)

	// wait for value to pass through buffers
	cache.Wait()

	// get value from cache
	value, found := cache.Get("key")
	if !found {
		panic("missing value")
	}
	fmt.Println(value)

	// del value from cache
	cache.Del("key")
}
*/
/*
import (
	"fmt"
	"time"

	"github.com/eko/gocache/lib/v4/cache"
	gocache_store "github.com/eko/gocache/store/go_cache/v4"
)

func f() {

		memcacheStore := memcache_store.NewMemcache(
			memcache.New("10.0.0.1:11211", "10.0.0.2:11211", "10.0.0.3:11212"),
			store.WithExpiration(10*time.Second),
		)

		cacheManager := cache.New[[]byte](memcacheStore)
		err := cacheManager.Set(ctx, "my-key", []byte("my-value"),
			store.WithExpiration(15*time.Second), // Override default value of 10 seconds defined in the store
		)
		if err != nil {
			panic(err)
		}

		value := cacheManager.Get(ctx, "my-key")

		cacheManager.Delete(ctx, "my-key")

		cacheManager.Clear(ctx) // Clears the entire cache, in case you want to flush all cache


	//inMemoryStore := store.NewGoCache(cache.New(5*time.Minute, 10*time.Minute), nil)

	gocacheClient := gocache_store.cache.New(5*time.Minute, 10*time.Minute)
	gocacheStore := gocache_store.NewGoCache(gocacheClient)

	cacheManager := cache.New[[]byte](gocacheStore)
	err := cacheManager.Set(ctx, "my-key", []byte("my-value"))
	if err != nil {
		panic(err)
	}

	value, err := cacheManager.Get(ctx, "my-key")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", value)
}
*/
