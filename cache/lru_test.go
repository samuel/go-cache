package cache

import (
	"testing"
)

func TestLRUCacheBasics(t *testing.T) {
	cache := NewLRUCache(2)
	if cache == nil {
		t.Fatal("Cache is nil")
	}

	BasicCacheTest(t, cache)
}

func TestLRUCacheExpiration(t *testing.T) {
	cache := NewLRUCache(2)
	if cache == nil {
		t.Fatal("Cache is nil")
	}

	if err := cache.Set("key1", "v1"); err != nil {
		t.Fatal("cache.Set returned err")
	}
	if err := cache.Set("key2", "v2"); err != nil {
		t.Fatal("cache.Set returned err")
	}

	if len(cache.index) != 2 || cache.lru.Len() != 2 {
		t.Fatal("index or lry length should be 2")
	}

	// Make sure set updates lru

	if err := cache.Set("key3", "v3"); err != nil {
		t.Fatal("cache.Set returned err")
	}

	if len(cache.index) != 2 || cache.lru.Len() != 2 {
		t.Fatal("index or lry length should be 2")
	}

	if val, err := cache.Get("key1"); err != nil {
		t.Fatal("cache.Get for existing key returned err")
	} else if val != nil {
		t.Fatal("didn't expire key1")
	}

	// Make sure get updates lru

	if val, err := cache.Get("key2"); err != nil {
		t.Fatal("cache.Get for existing key returned err")
	} else if val != "v2" {
		t.Fatal("expired key2 when shouldn't")
	}

	if err := cache.Set("key4", "v4"); err != nil {
		t.Fatal("cache.Set returned err")
	}

	if len(cache.index) != 2 || cache.lru.Len() != 2 {
		t.Fatal("index or lry length should be 2")
	}

	if val, err := cache.Get("key3"); err != nil {
		t.Fatal("cache.Get for existing key returned err")
	} else if val != nil {
		t.Fatal("didn't expire key3")
	}
}

func TestLRUCacheEvictionHook(t *testing.T) {
	callCount := 0

	cache := NewLRUCache(2)
	cache.SetEvictionHook(func(key, value interface{}) {
		callCount++
	})
	cache.Set("key1", "v1")
	cache.Set("key2", "v2")
	cache.Set("key3", "v3")
	if callCount != 1 {
		t.Fatalf("Eviction hook should be called")
	}
}

func TestLRUCacheKeys(t *testing.T) {
	cache := NewLRUCache(2)
	cache.Set("key1", "v1")
	cache.Set("key2", "v2")

	keys := cache.Keys()
	if length := len(keys); length != 2 {
		t.Fatalf("keys length should be 2 was %v", length)
	}
}
