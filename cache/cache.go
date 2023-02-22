package cache

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	lock sync.RWMutex
	data map[string][]byte
}

var _ Cacher = (*Cache)(nil)

func (c *Cache) Get(key []byte) ([]byte, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if v, ok := c.data[string(key)]; ok {
		fmt.Printf("get key: %s, value: %s success\n", key, v)
		return v, nil
	}
	return nil, errors.New("key not found")
}

func (c *Cache) Set(key []byte, value []byte, ttl time.Duration) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if ttl > 0 {
		go func() {
			<-time.After(ttl)
			c.Delete(key)
		}()
	}

	c.data[string(key)] = value
	fmt.Printf("set key: %s, value: %s, ttl: %d suceess\n", key, value, ttl)
	return nil
}

func (c *Cache) Delete(key []byte) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.data, string(key))
	fmt.Printf("delete overtime key: %s success\n", key)
	return nil
}

func (c *Cache) Has(key []byte) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if _, ok := c.data[string(key)]; ok {
		return true
	}
	return false
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string][]byte),
	}
}
