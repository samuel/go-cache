package cache

import (
	"math/rand"
	"strconv"
	"testing"
)

func TestLFUCacheBasics(t *testing.T) {
	cache := NewLFUCache(2)
	if cache == nil {
		t.Fatal("Cache is nil")
	}

	BasicCacheTest(t, cache)
}

func TestLFUCacheExpiration(t *testing.T) {
	cache := NewLFUCache(2)
	if cache == nil {
		t.Fatal("Cache is nil")
	}

	if err := cache.Set("key1", "v1"); err != nil {
		t.Fatal("cache.Set returned err")
	}
	if err := cache.Set("key2", "v2"); err != nil {
		t.Fatal("cache.Set returned err")
	}

	if len(cache.index) != 2 {
		t.Fatal("index length should be 2")
	}

	// Make sure set updates lfu

	if err := cache.Set("key3", "v3"); err != nil {
		t.Fatal("cache.Set returned err")
	}

	if len(cache.index) != 2 {
		t.Fatalf("index length should be 2 instead of %d", len(cache.index))
	}

	if val, err := cache.Get("key1"); err != nil {
		t.Fatal("cache.Get for existing key returned err")
	} else if val != nil {
		t.Fatal("didn't expire key1")
	}

	// Make sure get updates lfu

	if val, err := cache.Get("key2"); err != nil {
		t.Fatal("cache.Get for existing key returned err")
	} else if val != "v2" {
		t.Fatal("expired key2 when shouldn't")
	}

	if err := cache.Set("key4", "v4"); err != nil {
		t.Fatal("cache.Set returned err")
	}

	if len(cache.index) != 2 {
		t.Fatal("index or lry length should be 2")
	}

	if val, err := cache.Get("key3"); err != nil {
		t.Fatal("cache.Get for existing key returned err")
	} else if val != nil {
		t.Fatal("didn't expire key3")
	}
}

func TestLFUCacheEvictionHook(t *testing.T) {
	callCount := 0

	cache := NewLFUCache(2)
	cache.SetEvictionHook(func(key string, value interface{}) {
		callCount++
	})
	cache.Set("key1", "v1")
	cache.Set("key2", "v2")
	cache.Set("key3", "v3")
	if callCount != 1 {
		t.Fatalf("Eviction hook should be called")
	}
}

func BenchmarkLFUCacheGet10(b *testing.B) {
	benchmarkLFUCacheGet(b, 10)
}

func BenchmarkLFUCacheGet100(b *testing.B) {
	benchmarkLFUCacheGet(b, 1000)
}

func BenchmarkLFUCacheGet1000(b *testing.B) {
	benchmarkLFUCacheGet(b, 1000)
}

func BenchmarkLFUCacheGet10000(b *testing.B) {
	benchmarkLFUCacheGet(b, 10000)
}

func BenchmarkLFUCacheGet100000(b *testing.B) {
	benchmarkLFUCacheGet(b, 1000000)
}

func BenchmarkLFUCacheGet1000000(b *testing.B) {
	benchmarkLFUCacheGet(b, 1000000)
}

func benchmarkLFUCacheGet(b *testing.B, n int) {
	cache := NewLFUCache(n)
	keys := make([]string, n)
	for i := 0; i < n; i++ {
		k := strconv.Itoa(i)
		keys[i] = k
		cache.Set(k, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k := keys[rand.Intn(n)]
		cache.Get(k)
	}
	b.SetBytes(1)
}
