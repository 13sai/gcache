package lru

import (
	"container/list"
)

type Value interface {
	Len() int
}

type entry struct {
	key   string
	value Value
}

type Cache struct {
	maxBytes  uint32
	nBytes    uint32
	ll        *list.List
	cache     map[string]*list.Element
	OnEvicted func(key string, val Value)
}

func New(maxBytes uint32, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, found := c.cache[key]; found {
		c.ll.MoveToFront((ele))
		value, ok = ele.Value.(*entry).value, true
	}

	return
}

func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		e := ele.Value.(*entry)
		delete(c.cache, e.key)
		c.nBytes -= uint32(len(e.key) + e.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(e.key, e.value)
		}
	}
}

func (c *Cache) Add(key string, v Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		e := ele.Value.(*entry)
		c.nBytes += uint32(v.Len() - e.value.Len())
		e.value = v
	} else {
		ele := c.ll.PushFront(&entry{key, v})
		c.cache[key] = ele
		c.nBytes += uint32(len(key) + v.Len())
	}

	for c.maxBytes > 0 && c.maxBytes < c.nBytes {
		c.RemoveOldest()
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
