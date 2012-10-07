package cache

import (
	"container/list"
)

type LRUCache struct {
	maxItems int
	lru      *list.List
	index    map[string]*list.Element
	evictionHook func (string, interface{})
}

type keyValue struct {
	key   string
	value interface{}
}

func NewLRUCache(maxItems int) *LRUCache {
	cache := &LRUCache{
		maxItems: maxItems,
		lru:      list.New(),
		index:    make(map[string]*list.Element, maxItems),
	}
	return cache
}

func (c *LRUCache) SetEvictionHook(hook func(string, interface{})) {
	c.evictionHook = hook
}

func (c *LRUCache) Set(key string, value interface{}) error {
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

func (c *LRUCache) Get(key string) (interface{}, error) {
	el, ok := c.index[key]
	if !ok {
		return nil, nil
	}
	c.lru.MoveToBack(el)
	kv := el.Value.(*keyValue)
	return kv.value, nil
}

func (c *LRUCache) Delete(key string) error {
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