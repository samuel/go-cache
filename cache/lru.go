package cache

import (
	"container/list"
	"sync"
)

type LRUCache struct {
	maxItems     int
	lru          *list.List
	index        map[interface{}]*list.Element
	evictionHook func(interface{}, interface{})
	mu           sync.Mutex
}

type keyValue struct {
	key, value interface{}
}

func NewLRUCache(maxItems int) *LRUCache {
	cache := &LRUCache{
		maxItems: maxItems,
		lru:      list.New(),
		index:    make(map[interface{}]*list.Element, maxItems),
	}
	return cache
}

func (c *LRUCache) SetEvictionHook(hook func(key, value interface{})) {
	c.mu.Lock()
	c.evictionHook = hook
	c.mu.Unlock()
}

func (c *LRUCache) Set(key, value interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	el, ok := c.index[key]
	if ok {
		// Element exists so just move it to the back and update value
		kv := el.Value.(*keyValue)
		kv.value = value
		c.lru.MoveToBack(el)
	} else {
		if len(c.index) >= c.maxItems {
			// Cache is full so remove an existing key/value
			el := c.lru.Front()
			kv := el.Value.(*keyValue)
			c.delete(kv)
			// Reuse list element
			kv.key = key
			kv.value = value
			c.lru.MoveToBack(el)
			c.index[key] = el
		} else {
			// Cache is not full and this is a new key
			kv := &keyValue{key, value}
			el := c.lru.PushBack(kv)
			c.index[key] = el
		}
	}
	return nil
}

func (c *LRUCache) Get(key interface{}) (interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	el, ok := c.index[key]
	if !ok {
		return nil, nil
	}
	c.lru.MoveToBack(el)
	kv := el.Value.(*keyValue)
	return kv.value, nil
}

func (c *LRUCache) Delete(key interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	el, ok := c.index[key]
	if !ok {
		return nil
	}
	kv := c.lru.Remove(el).(*keyValue)
	c.delete(kv)
	return nil
}

func (c *LRUCache) delete(kv *keyValue) {
	if c.evictionHook != nil {
		c.evictionHook(kv.key, kv.value)
	}
	delete(c.index, kv.key)
}

func (c *LRUCache) Keys() []interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()
	keys := make([]interface{}, 0, len(c.index))
	for key, _ := range c.index {
		keys = append(keys, key)
	}
	return keys
}
