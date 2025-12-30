package grpc

import "C"
import (
	"api/config"
	"common/applog"
	"common/discovery"
	"common/tracer"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	userservice "grpc/user/user"
	"log"
)

var UserServiceClient userservice.UserClient

func InitRpcServiceClient() {
	ctx := context.Background()
	etcdRegister := discovery.NewResolver(config.GetConfig().Etcd.Addrs, applog.WrapGDPLogger(ctx))
	resolver.Register(etcdRegister)

	// 创建jaeger client
	otelHandler, err := tracer.JaegerClientHandler(
		config.GetConfig().Jaeger.Endpoints,
		config.GetConfig().Server.Name,
		config.GetConfig().Server.Env,
		config.GetConfig().Jaeger.IsOpenOnlySamplerError,
	)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.NewClient(
		"etcd:///user",
		grpc.WithStatsHandler(otelHandler),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	UserServiceClient = userservice.NewUserClient(conn)
	return
}
