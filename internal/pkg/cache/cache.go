package cache

import (
	"sync"
	"time"
)

const defaultCleanupInterval = 1 * time.Minute

type Config struct {
	defaultExpirationTime int16 `mapstructure:"default_exp_time"`
}

type Item struct {
	Value        any
	CreationTime time.Time
	Expiration   int64
}

type Cache struct {
	sync.RWMutex
	items             map[string]Item
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
}

func NewCache(config *Config) *Cache {
	items := make(map[string]Item)

	cache := &Cache{
		items:             items,
		defaultExpiration: time.Duration(config.defaultExpirationTime) * time.Second,
		cleanupInterval:   defaultCleanupInterval,
	}

	cache.cleanupInterval = defaultCleanupInterval
	go cache.StartGC()
	return cache
}

func (c *Cache) Set(k string, v any, expiration time.Duration) {
	if expiration == 0 {
		expiration = c.defaultExpiration
	}
	c.Lock()
	defer c.Unlock()

	c.items[k] = Item{
		Value:        v,
		CreationTime: time.Now(),
		Expiration:   time.Now().Add(expiration).UnixNano(),
	}
}

func (c *Cache) Get(k string) (any, bool) {
	c.RLock()
	defer c.RUnlock()

	item, ok := c.items[k]
	if !ok {
		return nil, false
	}

	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			return nil, false
		}
	}
	return item.Value, true
}

func (c *Cache) StartGC() {
	for {
		<-time.After(c.cleanupInterval)

		if c.items == nil {
			return
		}

		if keys := c.expiredKeys(); len(keys) != 0 {
			c.clearItems(keys)
		}
	}
}

func (c *Cache) expiredKeys() (keys []string) {
	c.RLock()

	defer c.RUnlock()

	for k, i := range c.items {
		if time.Now().UnixNano() > i.Expiration && i.Expiration > 0 {
			keys = append(keys, k)
		}
	}

	return
}

func (c *Cache) clearItems(keys []string) {
	c.Lock()

	defer c.Unlock()

	for _, k := range keys {
		delete(c.items, k)
	}
}
