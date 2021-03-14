package cache

import (
	"cache/iface"
	cache_proto "cache/proto"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
	cache_proto.RegisterCacheServer(RpcServer, a.CacheSrv)
	reflection.Register(RpcServer)
	err = RpcServer.Serve(server)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	select {}
}
