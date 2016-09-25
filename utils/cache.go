package utils

import (
	"errors"
	"sync"
	"time"
)

var ErrNotExist = errors.New("not exist")

const (
	defaultTTL = 30
)

type node struct {
	expiredAt time.Time
	val       string
}

func (n *node) expired() bool {
	return n.expiredAt.Before(time.Now())
}

type Cache struct {
	s    map[string]*node
	lock *sync.RWMutex
}

func NewCache() *Cache {
	c := &Cache{
		s:    make(map[string]*node),
		lock: &sync.RWMutex{},
	}
	go c.gc()

	return c
}

func (c *Cache) gc() {
	for {
		time.Sleep(60 * time.Second)
		for k, v := range c.s {
			if v.expired() {
				c.lock.Lock()
				delete(c.s, k)
				c.lock.Unlock()
			}
		}

	}
}

func (c *Cache) Get(key string) (string, error) {
	c.lock.RLock()
	node, ok := c.s[key]
	c.lock.RUnlock()

	if !ok || node.expired() {
		return "", ErrNotExist
	}

	return node.val, nil
}

func (c *Cache) Set(key, value string) error {
	c.lock.RLock()
	n, ok := c.s[key]
	c.lock.RUnlock()
	if !ok {
		n = &node{}
	}

	c.lock.Lock()
	n.val = value
	n.expiredAt = time.Now().Add(time.Duration(defaultTTL) * time.Second)
	c.s[key] = n
	c.lock.Unlock()

	return nil
}

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
