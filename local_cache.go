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

func (i *item) isExpired() bool {
	return time.Now().After(i.expiration)
}

func getItemDefaultOptions() *Options {
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

func NewCache(opts ...OptionsFn) Cache {
	o := getCacheDefaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	switch o.mode {
	case LRU:
		return newLRUCache(o)
	default:
		return newDefaultCache(o)
	}
}

type DefaultCache struct {
	valueMap map[string]item
	options  *Options
	rwMutex  sync.RWMutex
	sf       singleflight.Group
	ticker   *time.Ticker
}

func newDefaultCache(options *Options) Cache {

	c := &DefaultCache{
		valueMap: make(map[string]item),
		ticker:   time.NewTicker(5 * time.Second),
		options:  options,
	}

	go c.Clean()
	return c
}

func getCacheDefaultOptions() *Options {
	return &Options{
		expiration: 10 * time.Second,
		maxSum:     1024,
	}
}

func (d *DefaultCache) Set(key string, val interface{}, opts ...OptionsFn) {
	o := getCacheDefaultOptions()

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

func (d *DefaultCache) Clean() {
	var err error
	for {
		<-d.ticker.C
		err = d.defaultClean()
		fmt.Println(err)
	}
}

func (d *DefaultCache) defaultClean() error {

	for key, item := range d.valueMap {
		if item.isExpired() {
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
