package cache

type lruNode struct {
	key   string
	value interface{}
	next  *lruNode
	prev  *lruNode
}

type LRUCache struct {
	maxItems     int
	lruHead      *lruNode // least recently used
	lruTail      *lruNode // most recently used
	elements     []lruNode
	elementI     int
	index        map[string]*lruNode
	evictionHook func(string, interface{})
}

type keyValue struct {
	key   string
	value interface{}
}

func NewLRUCache(maxItems int) *LRUCache {
	cache := &LRUCache{
		maxItems: maxItems,
		lruHead:  nil,
		lruTail:  nil,
		elements: make([]lruNode, maxItems),
		elementI: 0,
		index:    make(map[string]*lruNode, maxItems),
	}
	return cache
}

func (c *LRUCache) SetEvictionHook(hook func(string, interface{})) {
	c.evictionHook = hook
}

func (c *LRUCache) Set(key string, value interface{}) error {
	el, ok := c.index[key]
	if ok {
		// Element exists so just move it to the front and update value
		el.value = value
		c.moveToFront(el)
	} else {
		if len(c.index) >= c.maxItems {
			// Cache is full so remove an existing key/value
			el := c.lruTail
			c.lruTail = el.prev
			if el.prev != nil {
				el.prev.next = nil
			} else {
				c.lruHead = nil
			}
			c.delete(el)
			// Reuse list element
			el.key = key
			el.value = value
			el.prev = nil
			el.next = c.lruHead
			if c.lruHead != nil {
				c.lruHead.prev = el
			}
			c.lruHead = el
			c.index[key] = el
		} else {
			// Cache is not full and this is a new key
			el := &c.elements[c.elementI]
			c.elementI++
			el.key = key
			el.value = value
			el.next = c.lruHead
			el.prev = nil
			if c.lruHead != nil {
				c.lruHead.prev = el
			} else {
				c.lruTail = el
			}
			c.lruHead = el
			c.index[key] = el
		}
	}
	return nil
}

func (c *LRUCache) lruLen() int {
	i := 0
	for n := c.lruHead; n != nil; n = n.next {
		i++
	}
	return i
}

func (c *LRUCache) moveToFront(el *lruNode) {
	if el.prev != nil {
		el.prev.next = el.next
		if el.next != nil {
			el.next.prev = el.prev
		} else {
			c.lruTail = el.prev
		}
		el.next = c.lruHead
		el.prev = nil
		c.lruHead.prev = el
		c.lruHead = el
	}
}

func (c *LRUCache) Get(key string) (interface{}, error) {
	el, ok := c.index[key]
	if !ok {
		return nil, nil
	}
	c.moveToFront(el)
	return el.value, nil
}

func (c *LRUCache) Delete(key string) error {
	el, ok := c.index[key]
	if !ok {
		return nil
	}
	if el.prev == nil {
		c.lruHead = el.next
	} else {
		el.prev.next = el.next
	}
	if el.next == nil {
		c.lruTail = el.prev
	} else {
		el.next.prev = el.prev
	}
	c.delete(el)
	return nil
}

func (c *LRUCache) delete(node *lruNode) {
	if c.evictionHook != nil {
		c.evictionHook(node.key, node.value)
	}
	delete(c.index, node.key)
}

func (c *LRUCache) Keys() []string {
	keys := make([]string, 0, len(c.index))
	for key, _ := range c.index {
		keys = append(keys, key)
	}
	return keys
}
