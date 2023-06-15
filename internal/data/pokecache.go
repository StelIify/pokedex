package data

import (
	"sync"
	"time"
)

type CacheEntry struct {
	createdAt time.Time
	value     []byte
}

type Cache struct {
	mu   sync.RWMutex
	data map[string]CacheEntry
}

func NewCache(duration time.Duration) *Cache {
	cache := &Cache{
		data: make(map[string]CacheEntry),
	}
	go cache.ReapLoop(duration)
	return cache
}

func (c *Cache) Put(key string, value []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = CacheEntry{
		createdAt: time.Now(),
		value:     value,
	}

	return nil
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cacheEntry, ok := c.data[key]
	if !ok {
		return cacheEntry.value, false
	}
	return cacheEntry.value, true
}
func (c *Cache) ReapLoop(duration time.Duration) {
	ticker := time.NewTicker(duration)
	for range ticker.C {
		c.Reap(duration)
	}
}
func (c *Cache) Reap(duration time.Duration) {
	cacheDuration := time.Now().Add(-duration)
	for k, v := range c.data {
		if v.createdAt.Before(cacheDuration) {
			delete(c.data, k)
		}
	}
}
