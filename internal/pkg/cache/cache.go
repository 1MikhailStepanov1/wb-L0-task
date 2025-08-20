package cache

import (
	"sync"
	"time"
)

const defaultCleanupInterval = 1 * time.Minute

type Config struct {
	defaultExpirationTime int16 `mapstructure:"default_exp_time"`
}

type Item[T any] struct {
	Value        T
	CreationTime time.Time
	Expiration   int64
}

type Cache[T any] struct {
	sync.RWMutex
	items             map[string]Item[T]
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
}

func NewCache[T any](config *Config) *Cache[T] {
	items := make(map[string]Item[T])

	cache := &Cache[T]{
		items:             items,
		defaultExpiration: time.Duration(config.defaultExpirationTime) * time.Second,
		cleanupInterval:   defaultCleanupInterval,
	}

	cache.cleanupInterval = defaultCleanupInterval
	go cache.StartGC()
	return cache
}

// Set Add element to cache.
// If expiration == 0 default expiration time will be used.
func (c *Cache[T]) Set(k string, v T, expiration time.Duration) {
	if expiration == 0 {
		expiration = c.defaultExpiration
	}
	c.Lock()
	defer c.Unlock()

	c.items[k] = Item[T]{
		Value:        v,
		CreationTime: time.Now(),
		Expiration:   time.Now().Add(expiration).UnixNano(),
	}
}

func (c *Cache[T]) Get(k string) (T, bool) { //nolint:ireturn
	c.RLock()
	defer c.RUnlock()

	var zeroValue T
	item, ok := c.items[k]
	if !ok {
		return zeroValue, false
	}

	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			return zeroValue, false
		}
	}
	return item.Value, true
}

func (c *Cache[T]) StartGC() {
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

func (c *Cache[T]) expiredKeys() (keys []string) {
	c.RLock()

	defer c.RUnlock()

	for k, i := range c.items {
		if time.Now().UnixNano() > i.Expiration && i.Expiration > 0 {
			keys = append(keys, k)
		}
	}
	return
}

func (c *Cache[T]) clearItems(keys []string) {
	c.Lock()

	defer c.Unlock()

	for _, k := range keys {
		delete(c.items, k)
	}
}
