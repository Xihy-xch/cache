package iface

import (
	"bytes"
	"context"
	"encoding/gob"
	"local-cache/local_cache"
	cache "local-cache/proto"
)

type CacheSrv struct {
	cache local_cache.Cache
}

func NewCacheSrv(cache *local_cache.LRUCache) *CacheSrv {
	return &CacheSrv{cache: cache}
}

func (c *CacheSrv) Get(ctx context.Context, request *cache.CacheGetRequest) (*cache.CacheGetResponse, error) {
	val, err := c.cache.Get(request.GetKey())
	if err != nil {
		return nil, err
	}
	res, err := getBytes(val)
	if err != nil {
		return nil, err
	}

	return &cache.CacheGetResponse{Value: res}, nil
}

func (c *CacheSrv) Set(ctx context.Context, request *cache.CacheSetRequest) (*cache.CacheSetResponse, error) {
	panic("implement me")
}

func getBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
