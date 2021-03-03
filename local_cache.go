package local_cache

import (
	"errors"
	"fmt"
	"golang.org/x/sync/singleflight"
	"sync"
	"time"
)

type item struct {
	value      interface{}
	expiration time.Time
}

func (i *item) getDefaultOptions() *Options {
	return &Options{
		expiration: 10 * time.Second,
	}
}

type Cache struct {
	valueMap map[string]item
	options  *Options
	rwMutex  *sync.RWMutex
	sf       *singleflight.Group
	ticker   *time.Ticker
}

func NewCache(opts ...OptionsFn) *Cache {
	c := &Cache{
		valueMap: make(map[string]item),
		rwMutex:  new(sync.RWMutex),
		sf:       new(singleflight.Group),
		ticker:   time.NewTicker(5 * time.Second),
	}

	o := c.getDefaultOptions()
	for _, opt := range opts {
		opt(o)
	}
	c.options = o
	go c.clean()
	return c
}

func (c *Cache) getDefaultOptions() *Options {
	return &Options{
		expiration: 10 * time.Second,
		maxSum:     1024,
	}
}

func (c *Cache) Set(key string, val interface{}, opts ...OptionsFn) {
	o := c.getDefaultOptions()

	for _, opt := range opts {
		opt(o)
	}

	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()
	c.valueMap[key] = item{
		value:      val,
		expiration: time.Now().Add(o.GetExpiration()),
	}
}

func (c *Cache) Get(key string) (interface{}, error) {
	val, err, _ := c.sf.Do(key, func() (interface{}, error) {
		c.rwMutex.RLock()
		defer c.rwMutex.RUnlock()
		if val, exist := c.valueMap[key]; exist {
			return val, nil
		}

		return nil, errors.New("该key不存在")
	})

	if res, ok := val.(item); ok {
		return res.value, nil
	}
	return nil, err
}

// @todo error处理
func (c *Cache) clean() error {
	for {
		<-c.ticker.C
		fmt.Println("开始清理")
		switch c.options.cleanMode {
		case Default:

			c.defaultClean()
		case LRU:
			c.lruClean()
		default:
			c.defaultClean()
		}
	}
}

func (c *Cache) defaultClean() error {
	now := time.Now()

	for key, item := range c.valueMap {
		if now.After(item.expiration) {
			delete(c.valueMap, key)
		}
	}

	return nil
}

func (c *Cache) lruClean() error {
	return nil
}
