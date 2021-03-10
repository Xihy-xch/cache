//+build wireinject

package inject

import (
	"github.com/google/wire"
	local_cache "local-cache"
	"local-cache/iface"
	local_cache2 "local-cache/local_cache"
)

func InitApp(maxSum int64) (app *local_cache.App, cleanup func(), err error) {
	wire.Build(
		local_cache2.NewLRUCache,
		iface.NewCacheSrv,
		local_cache.NewApp,
	)
	return &local_cache.App{}, nil, nil
}
