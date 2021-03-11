//+build wireinject

package inject

import (
	"cache"
	"cache/iface"
	"cache/local_cache"
	"github.com/google/wire"
)

func InitApp(maxSum int64) (app *cache.App, cleanup func(), err error) {
	wire.Build(
		local_cache.NewLRUCache,
		iface.NewCacheSrv,
		cache.NewApp,
	)
	return &cache.App{}, nil, nil
}
