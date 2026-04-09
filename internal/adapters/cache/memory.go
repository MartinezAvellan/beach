package cache

import (
	"sync"
	"time"
)

type entry struct {
	value     any
	expiresAt time.Time
}

type MemoryCache struct {
	mu      sync.RWMutex
	items   map[string]entry
	ttl     time.Duration
}

func New(ttl time.Duration) *MemoryCache {
	return &MemoryCache{
		items: make(map[string]entry),
		ttl:   ttl,
	}
}

func (c *MemoryCache) Get(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	e, ok := c.items[key]
	if !ok || time.Now().After(e.expiresAt) {
		return nil, false
	}
	return e.value, true
}

func (c *MemoryCache) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = entry{
		value:     value,
		expiresAt: time.Now().Add(c.ttl),
	}
}
