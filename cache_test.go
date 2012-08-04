package cache

import (
	"testing"
)

func BasicCacheTest(t *testing.T, cache Cache) {
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

	/*if cache.Count() != 2) {
		t.Fatal("count should be 2")
	}*/

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

	/*if cache.Count() != 1 {
	    t.Fatal("count should be 1")
	}*/
}
