package iface

import (
	"cache/local_cache"
	cacheProto "cache/proto"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
)

type CacheSrv struct {
	cache local_cache.Cache
}

func NewCacheSrv(cache *local_cache.LRUCache) *CacheSrv {
	return &CacheSrv{cache: cache}
}

func (c *CacheSrv) Get(ctx context.Context, request *cacheProto.CacheGetRequest) (*cacheProto.CacheGetResponse, error) {
	val, err := c.cache.Get(request.GetKey())
	if err != nil {
		return nil, err
	}
	res, err := marshal(val)
	if err != nil {
		return nil, err
	}

	return &cacheProto.CacheGetResponse{Value: res}, nil
}

func (c *CacheSrv) Set(ctx context.Context, request *cacheProto.CacheSetRequest) (*cacheProto.CacheSetResponse, error) {
	var val interface{}
	err := unmarshal(request.GetValue(), &val)
	if err != nil {
		return nil, err
	}

	c.cache.Set(request.GetKey(), val)

	return &cacheProto.CacheSetResponse{}, nil
}

func marshal(v interface{}) ([]byte, error) {
	b, err := json.Marshal(v)
	return b, errors.WithStack(err)
}

func unmarshal(data []byte, v interface{}) error {
	return errors.WithStack(json.Unmarshal(data, v))
}
