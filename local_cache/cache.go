package local_cache

import (
	"golang.org/x/sync/singleflight"
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