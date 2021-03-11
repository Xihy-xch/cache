package local_cache

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"local-cache/iface"
	cache "local-cache/proto"
	"log"
	"net"
)

type App struct {
	CacheSrv *iface.CacheSrv
}

func NewApp(cacheSrv *iface.CacheSrv) *App {
	return &App{CacheSrv: cacheSrv}
}

func (a *App) Start() {
	rpcAddr := "127.0.0.1:8601"
	server, err := net.Listen("tcp", rpcAddr)
	if err != nil {
		fmt.Println("failed to listen", rpcAddr)
		panic(err)
	}

	// 建立rpc server
	var RpcServer = grpc.NewServer()
	cache.RegisterCacheServer(RpcServer, a.CacheSrv)
	reflection.Register(RpcServer)
	err = RpcServer.Serve(server)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	select {}
}
