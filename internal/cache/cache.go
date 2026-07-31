package cache

import (
	"container/list"
	"sync"
	"time"

	stats "github.com/Abdullah-Builds/cache/internal/statistics"
)

type Entry struct {
	Key  string
	Item Item
}

type Cache struct {
	data       map[string]*list.Element
	mu         sync.RWMutex
	maxKeys    int
	lru        *list.List
	stats      stats.Stats
	maxMemory  int64
	usedMemory int64
}

func New() *Cache {
	return NewWithMaxKeys(100)
}

func NewWithMaxKeys(maxKeys int) *Cache {
	if maxKeys <= 0 {
		maxKeys = 100
	}
	return &Cache{data: map[string]*list.Element{}, maxKeys: maxKeys, stats: stats.Stats{}, lru: list.New()}
}

func (c *Cache) Set(key, value string, ttl time.Duration) string {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.stats.Requests.Add(1)

	if existing, ok := c.data[key]; ok {
		entry := existing.Value.(*Entry)
		c.usedMemory -= int64(len(entry.Key) + len(entry.Item.Value))
		c.lru.Remove(existing)
		delete(c.data, key)
	} else if len(c.data) >= c.maxKeys {
		oldest := c.lru.Back()
		entry := oldest.Value.(*Entry)
		delete(c.data, entry.Key)
		c.usedMemory -= int64(len(entry.Key) + len(entry.Item.Value))
		c.lru.Remove(oldest)
	}

	entry := &Entry{
		Key: key,
		Item: Item{
			Value: value,
		},
	}

	if ttl > 0 {
		entry.Item.ExpiresAt = time.Now().Add(ttl)
	}

	c.stats.Sets.Add(1)

	size := int64(len(key) + len(value))

	c.usedMemory += size
	element := c.lru.PushFront(entry)
	c.data[key] = element

	return "OK"
}

func (c *Cache) Get(key string) (string, bool) {

	c.stats.Requests.Add(1)

	c.mu.Lock()
	defer c.mu.Unlock()

	element, ok := c.data[key]
	if !ok {
		c.stats.Misses.Add(1)
		return "", false
	}
	entry := element.Value.(*Entry)
	if !entry.Item.ExpiresAt.IsZero() && !time.Now().Before(entry.Item.ExpiresAt) {
		c.lru.Remove(element)
		delete(c.data, key)
		c.usedMemory -= int64(len(key) + len(entry.Item.Value))
		c.stats.Misses.Add(1)
		c.stats.Deletes.Add(1)
		return "", false
	}
	c.lru.MoveToFront(element)
	c.stats.Hits.Add(1)
	return entry.Item.Value, true
}

func (c *Cache) Delete(key string) {

	c.mu.Lock()
	defer c.mu.Unlock()

	element, ok := c.data[key]
	if !ok {
		return
	}
	c.lru.Remove(element)
	c.usedMemory -= int64(len(key) + len(element.Value.(*Entry).Item.Value))

	delete(c.data, key)
	c.stats.Deletes.Add(1)
}

// Now expired keys disappear automatically.
func (c *Cache) StartCleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		defer ticker.Stop()

		for range ticker.C {
			c.cleanupExpired()
		}
	}()
}

func (c *Cache) cleanupExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()

	for key, element := range c.data {
		if element == nil || element.Value == nil {
			delete(c.data, key)
			continue
		}

		entry := element.Value.(*Entry)
		if !entry.Item.ExpiresAt.IsZero() && now.After(entry.Item.ExpiresAt) {
			c.lru.Remove(element)
			c.usedMemory -= int64(len(key) + len(entry.Item.Value))
			delete(c.data, key)
			c.stats.Deletes.Add(1)
		}
	}
}

func (c *Cache) Stats() *stats.Stats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return &c.stats
}
