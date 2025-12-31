package grpc

import (
	"context"
	"log"

	"common/applog"
	"common/discovery"
	"common/tracer"
	"user/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

func InitRpcService() {
	ctx := context.Background()
	etcdRegister := discovery.NewResolver(config.GetConfig().Etcd.Addrs, applog.WrapGDPLogger(ctx))
	resolver.Register(etcdRegister)

	otelHandler, err := tracer.JaegerClientHandler(
		config.GetConfig().Jaeger.Endpoints,
		config.GetConfig().Server.Name,
		config.GetConfig().Server.Env,
		config.GetConfig().Jaeger.IsOpenOnlySamplerError,
	)
	if err != nil {
		log.Fatal(err)
	}

	_, err = grpc.NewClient(
		discovery.BuildResolverUrl("project"),
		grpc.WithStatsHandler(otelHandler),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	//ProjectServiceClient = projectservice.NewProjectClient(conn)
}
