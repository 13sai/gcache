package gcache

import (
	"sync"

	"github.com/13sai/gcache/lru"
)

type cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes uint32
}

func (c *cache) add(k string, v ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(k, v)
}

func (c *cache) get(k string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(k); ok {
		return v.(ByteView), ok
	}
	return
}
