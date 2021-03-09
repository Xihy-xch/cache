package local_cache

import (
	"fmt"
	"golang.org/x/sync/singleflight"
	"sync"
	"time"
)

var (
	sf singleflight.Group
)

type item struct {
	value      interface{}
	expiration time.Time
}

func (i *item) isExpired() bool {
	return time.Now().After(i.expiration)
}

type Cache interface {
	Get(key string, opts ...OptionsFn) (interface{}, error)
	Set(key string, val interface{}, opts ...OptionsFn)
	Delete(key string)
	Clean()
	Close()
}

func NewCache(cache Cache) Cache {
	return cache
}

type DefaultCache struct {
	valueMap    map[string]item
	rwMutex     sync.RWMutex
	cleanTicker *time.Ticker
	stop        chan int
}

func NewDefaultCache(cleanTicker *time.Ticker) *DefaultCache {

	c := &DefaultCache{
		valueMap:    make(map[string]item),
		cleanTicker: cleanTicker,
		stop:        make(chan int),
	}

	go c.Clean()
	return c
}

func (d *DefaultCache) Set(key string, val interface{}, opts ...OptionsFn) {
	o := getDefaultOptions()

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

func (d *DefaultCache) Get(key string, opts ...OptionsFn) (interface{}, error) {
	o := getDefaultOptions()

	for _, opt := range opts {
		opt(o)
	}

	val, err := d.doGet(key)
	if err == nil {
		return val, nil
	}

	if o.getter == nil {
		return nil, err
	}

	val, err, _ = sf.Do(key, func() (interface{}, error) {
		val, err := o.getter.Get(key)
		return val, err
	})
	if err != nil {
		return nil, err
	}

	return val, nil
}

func (d *DefaultCache) doGet(key string) (interface{}, error) {
	d.rwMutex.RLock()
	defer d.rwMutex.RUnlock()
	var (
		val   item
		exist bool
	)
	if val, exist = d.valueMap[key]; !exist {
		return nil, ErrKeyNotExist
	}
	if val.isExpired() {
		return nil, ErrKeyExpired
	}

	return val.value, nil
}

func (d *DefaultCache) Delete(key string) {
	d.rwMutex.Lock()
	defer d.rwMutex.Unlock()
	delete(d.valueMap, key)
}

func (d *DefaultCache) Clean() {
	for {
		select {
		case <-d.stop:
			return
		case <-d.cleanTicker.C:
			d.defaultClean()
		}
	}
}

func (d *DefaultCache) defaultClean() {
	fmt.Println("开始清理")
	for key, item := range d.valueMap {
		if item.isExpired() {
			delete(d.valueMap, key)
		}
	}
}

func (d *DefaultCache) Close() {
	d.cleanTicker.Stop()
	d.stop <- 1
	close(d.stop)
}
