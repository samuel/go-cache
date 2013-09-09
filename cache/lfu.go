package cache

import (
	"container/list"
	"sync"
)

type LFUCache struct {
	maxItems     int
	lfu          *list.List
	index        map[interface{}]*list.Element
	evictionHook func(interface{}, interface{})
	mu           sync.Mutex
}

type lfuNode struct {
	count int
	items *list.List
}

type lfuItem struct {
	parent     *list.Element
	key, value interface{}
}

func NewLFUCache(maxItems int) *LFUCache {
	c := &LFUCache{
		maxItems: maxItems,
		lfu:      list.New(),
		index:    make(map[interface{}]*list.Element, maxItems),
	}
	firstNode := &lfuNode{1, list.New()}
	c.lfu.PushFront(firstNode)
	return c
}

func (c *LFUCache) SetEvictionHook(hook func(key, value interface{})) {
	c.mu.Lock()
	c.evictionHook = hook
	c.mu.Unlock()
}

func (c *LFUCache) Set(key, value interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	el, ok := c.index[key]
	if ok {
		el.Value.(*lfuItem).value = value
		c.incrLfu(el)
	} else {
		if len(c.index) >= c.maxItems {
			c.expireOneItem()
		}
		n := c.lfu.Front()
		item := &lfuItem{n, key, value}
		c.index[key] = n.Value.(*lfuNode).items.PushFront(item)
	}
	return nil
}

func (c *LFUCache) Get(key interface{}) (interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	el, ok := c.index[key]
	if ok {
		c.incrLfu(el)
		return el.Value.(*lfuItem).value, nil
	}
	return nil, nil
}

func (c *LFUCache) Delete(key interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	el, ok := c.index[key]
	if ok {
		c.removeElement(el)
		c.delete(el.Value.(*lfuItem))
	}
	return nil
}

func (c *LFUCache) incrLfu(el *list.Element) {
	item := el.Value.(*lfuItem)
	lfu := item.parent.Value.(*lfuNode)
	newCount := lfu.count + 1
	nextLfuEl := item.parent.Next()
	if nextLfuEl == nil || nextLfuEl.Value.(*lfuNode).count > newCount {
		if lfu.count != 1 && lfu.items.Len() == 1 {
			// OPTIMIZATION: If this is the only element in the LFU node and
			// it's not the first node then all we need to do is increment the count.
			lfu.count++
		} else {
			// Create new LFU node and move element to it
			nextNode := &lfuNode{newCount, list.New()}
			c.index[item.key] = nextNode.items.PushFront(item)
			c.lfu.InsertAfter(nextNode, item.parent)
			c.removeElement(el)
		}
	} else {
		// Move element to next LFU node
		c.index[item.key] = nextLfuEl.Value.(*lfuNode).items.PushFront(item)
		c.removeElement(el)
	}
}

func (c *LFUCache) removeElement(el *list.Element) {
	lfuEl := el.Value.(*lfuItem).parent
	lfu := lfuEl.Value.(*lfuNode)
	lfu.items.Remove(el)
	if lfu.count != 1 && lfu.items.Len() == 0 {
		c.lfu.Remove(lfuEl)
	}
}

// Expire the oldest item with the lowest count available
func (c *LFUCache) expireOneItem() {
	// This will only ever loop once or twice since count 1 is
	// guaranteed to exist but may be empty. The next node
	// must always be non-empty.
	for lfuEl := c.lfu.Front(); lfuEl != nil; lfuEl = lfuEl.Next() {
		lfuNode := lfuEl.Value.(*lfuNode)
		el := lfuNode.items.Back()
		if el != nil {
			item := el.Value.(*lfuItem)
			lfuNode.items.Remove(el)
			c.delete(item)
			return
		}
	}
}

func (c *LFUCache) delete(item *lfuItem) {
	if c.evictionHook != nil {
		c.evictionHook(item.key, item.value)
	}
	delete(c.index, item.key)
}
