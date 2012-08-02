package cache

import (
	"testing"
)

func TestLRUCacheBasics(t *testing.T) {
	cache := NewLRUCache(2)
	if cache == nil {
		t.Fatal("Cache is nil")
	}

	if val, err := cache.Get("nonexistant"); err != nil {
		t.Fatal("non-existant key returned error")
	} else if val != nil {
		t.Fatal("non-existant key returned value")
	}

	if err := cache.Set("exists", "here"); err != nil {
		t.Fatal("cache.Set returned err")
	}
	if err := cache.Set("foo", "bar"); err != nil {
		t.Fatal("cache.Set returned err")
	}

	if len(cache.index) != 2 || cache.lru.Len() != 2 {
		t.Fatal("index or lry length should be 2")
	}

	if val, err := cache.Get("exists"); err != nil {
		t.Fatal("cache.Get for existing key returned err")
	} else if val == nil {
		t.Fatal("cache.Get for existing key returned nil")
	} else if val.(string) != "here" {
		t.Fatal("cache.Get returned wrong value")
	}

	// Delete

	if err := cache.Delete("exists"); err != nil {
		t.Fatal("cache.Delete returned err")
	}
	if val, err := cache.Get("exists"); err != nil {
		t.Fatal("deleted key returned error")
	} else if val != nil {
		t.Fatal("deleted key returned value")
	}
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
