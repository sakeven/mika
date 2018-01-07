package utils

import (
	"errors"
	"sync"
	"time"
)

// ErrNotExist means key doesn't exist in the cache.
var ErrNotExist = errors.New("not exist")

const (
	defaultTTL = 30 * time.Second
)

type node struct {
	expiredAt time.Time
	val       string
}

func (n *node) expired() bool {
	return n.expiredAt.Before(time.Now())
}

// Cache stores k-v pair in memory.
type Cache struct {
	s          map[string]*node
	gcInterval time.Duration
	lock       *sync.RWMutex
}

// NewCache creates a new cache instance.
func NewCache() *Cache {
	c := &Cache{
		s:          make(map[string]*node),
		gcInterval: 60 * time.Second,
		lock:       &sync.RWMutex{},
	}
	go c.gc()

	return c
}

func (c *Cache) gc() {
	for {
		time.Sleep(c.gcInterval)
		for k, v := range c.s {
			if v.expired() {
				c.lock.Lock()
				delete(c.s, k)
				c.lock.Unlock()
			}
		}
	}
}

// Get fetches a value from cache for specific key.
func (c *Cache) Get(key string) (string, error) {
	c.lock.RLock()
	node, ok := c.s[key]
	c.lock.RUnlock()

	if !ok || node.expired() {
		return "", ErrNotExist
	}

	return node.val, nil
}

// Set stores <key, value> to cache.
func (c *Cache) Set(key, value string) error {
	c.lock.RLock()
	n, ok := c.s[key]
	c.lock.RUnlock()
	if !ok {
		n = &node{}
	}

	c.lock.Lock()
	n.val = value
	n.expiredAt = time.Now().Add(defaultTTL)
	c.s[key] = n
	c.lock.Unlock()

	return nil
}

// SetWithTTL stores <key, value> pair to cache, and after ttl seconds, the pair will be removed.
func (c *Cache) SetWithTTL(key, value string, ttl int64) error {
	c.lock.RLock()
	n, ok := c.s[key]
	c.lock.RUnlock()
	if !ok {
		n = &node{}
	}

	c.lock.Lock()
	n.val = value
	n.expiredAt = time.Now().Add(time.Duration(ttl) * time.Second)
	c.s[key] = n
	c.lock.Unlock()

	return nil
}
