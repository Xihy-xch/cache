package local_cache

import (
	"errors"
	"fmt"
	lru "github.com/hashicorp/golang-lru"
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

type Cache interface {
	Get(key string) (interface{}, error)
	Set(key string, val interface{}, opts ...OptionsFn)
	Delete(key string)
	Clean()
}

type DefaultCache struct {
	valueMap map[string]item
	options  *Options
	rwMutex  *sync.RWMutex
	sf       *singleflight.Group
	ticker   *time.Ticker
}

func NewCache(opts ...OptionsFn) *DefaultCache {
	c := &DefaultCache{
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
	go c.Clean()
	return c
}

func (d *DefaultCache) getDefaultOptions() *Options {
	return &Options{
		expiration: 10 * time.Second,
		maxSum:     1024,
	}
}

func (d *DefaultCache) Set(key string, val interface{}, opts ...OptionsFn) {
	o := d.getDefaultOptions()

	for _, opt := range opts {
		opt(o)
	}

	d.rwMutex.Lock()
	defer d.rwMutex.Unlock()
	d.valueMap[key] = item{
		value:      val,
		expiration: time.Now().Add(o.GetExpiration()),
	}
}

func (d *DefaultCache) Get(key string) (interface{}, error) {
	val, err, _ := d.sf.Do(key, func() (interface{}, error) {
		d.rwMutex.RLock()
		defer d.rwMutex.RUnlock()
		if val, exist := d.valueMap[key]; exist {
			return val, nil
		}

		return nil, errors.New("该key不存在")
	})

	if res, ok := val.(item); ok {
		return res.value, nil
	}
	return nil, err
}

func (d *DefaultCache) Delete(key string) {
	d.rwMutex.Lock()
	defer d.rwMutex.Unlock()
	delete(d.valueMap, key)
}

// @todo error处理
func (d *DefaultCache) Clean() {
	var err error
	for {
		<-d.ticker.C
		fmt.Println("开始清理")
		switch d.options.cleanMode {
		case Default:
			err = d.defaultClean()
		case LRU:
			err = d.lruClean()
		default:
			err = d.defaultClean()
		}
		fmt.Println(err)
	}
}

func (d *DefaultCache) defaultClean() error {
	now := time.Now()

	for key, item := range d.valueMap {
		if now.After(item.expiration) {
			delete(d.valueMap, key)
		}
	}

	return nil
}

func (d *DefaultCache) lruClean() error {
	l, _ := lru.New(128)
	for i := 0; i < 256; i++ {
		l.Add(i, nil)
	}
	if l.Len() != 128 {
		panic(fmt.Sprintf("bad len: %v", l.Len()))
	}
	return nil
}
