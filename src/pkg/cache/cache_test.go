package cache

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type type1 struct {
	name    string
	age     int
	names   []string
	numbers []int
}

type type2 struct {
	names   []string
	numbers []int
	t1      type1
}

var (
	t1 = type1{name: "Alice", age: 30, names: []string{"Bob", "Charlie"}, numbers: []int{1, 2, 3}}

	t21 = type2{names: []string{"Alice", "Bob"}, numbers: []int{1, 2, 3}, t1: type1{name: "Charlie", age: 30, names: []string{"X", "Y"}, numbers: []int{4, 5, 6}}}
	t22 = type2{names: []string{"Alice", "Bob"}, numbers: []int{1, 2, 3}, t1: type1{name: "Charlie", age: 30, names: []string{"X", "Y"}, numbers: []int{4, 5, 6}}}
	t23 = type2{names: []string{"Alice", "Bob"}, numbers: []int{4, 5, 6}, t1: type1{name: "Charlie", age: 30, names: []string{"X", "Y"}, numbers: []int{4, 5, 6}}}
	t24 = type2{names: []string{"Alice", "Bob"}, numbers: []int{1, 2, 3}, t1: type1{name: "Charlie", age: 31, names: []string{"X", "Y"}, numbers: []int{4, 5, 6}}}

	t31 = &type2{names: []string{"Alice", "Bob"}, numbers: []int{1, 2, 3}, t1: type1{name: "Charlie", age: 30, names: []string{"X", "Y"}, numbers: []int{4, 5, 6}}}
	t32 = &type2{names: []string{"Alice", "Bob"}, numbers: []int{1, 2, 3}, t1: type1{name: "Charlie", age: 30, names: []string{"X", "Y"}, numbers: []int{4, 5, 6}}}
	t33 = &type2{names: []string{"Alice", "Bob"}, numbers: []int{4, 5, 6}, t1: type1{name: "Charlie", age: 30, names: []string{"X", "Y"}, numbers: []int{4, 5, 6}}}
	t34 = &type2{names: []string{"Alice", "Bob"}, numbers: []int{1, 2, 3}, t1: type1{name: "Charlie", age: 31, names: []string{"X", "Y"}, numbers: []int{4, 5, 6}}}
)

func waitConsolidation(cache CacheService) {
	if lowlevel, ok := cache.(*cacheServiceImpl); ok {
		lowlevel.provider.Wait()
	} else {
		time.Sleep(100 * time.Millisecond)
	}
}

func TestValues(t *testing.T) {
	cache := GetCache()
	cache.Set("t1", t1)
	waitConsolidation(cache)
	cachedt1, found := cache.Get("t1")
	assert.True(t, found)
	assert.True(t, reflect.DeepEqual(t1, cachedt1))
	cache.Del("t1")
	waitConsolidation(cache)
	_, found = cache.Get("t1")
	assert.True(t, !found)
}

func TestPointers1(t *testing.T) {
	cache := GetCache()
	cache.Set("type2.t21", &t21)
	waitConsolidation(cache)
	cachedt21, found := cache.Get("type2.t21")
	assert.True(t, found)
	assert.True(t, reflect.DeepEqual(&t21, cachedt21))
	assert.True(t, reflect.DeepEqual(&t22, cachedt21))
	assert.False(t, reflect.DeepEqual(&t23, cachedt21))
	assert.False(t, reflect.DeepEqual(&t24, cachedt21))
	t24.t1.age = 30
	assert.True(t, reflect.DeepEqual(&t24, cachedt21))
	cache.Del("type2.t21")
	waitConsolidation(cache)
	_, found = cache.Get("type2.t21")
	assert.True(t, !found)
}

func TestPointers2(t *testing.T) {
	cache := GetCache()
	cache.Set("type2.t31", t31)
	waitConsolidation(cache)
	cachedt31, found := cache.Get("type2.t31")
	assert.True(t, found)
	assert.True(t, reflect.DeepEqual(t31, cachedt31))
	assert.True(t, reflect.DeepEqual(t32, cachedt31))
	assert.False(t, reflect.DeepEqual(t33, cachedt31))
	assert.False(t, reflect.DeepEqual(t34, cachedt31))
	t34.t1.age = 30
	assert.True(t, reflect.DeepEqual(t34, cachedt31))
	cache.Del("type2.t31")
	waitConsolidation(cache)
	_, found = cache.Get("type2.t31")
	assert.True(t, !found)
}
