package main

import (
	srv "common"
	"context"
	"user/pkg/grpc"
	"user/pkg/initialize"
)

func main() {
	ctx := context.Background()
	initialize.MustInit(ctx)
	//grpc服务注册
	gc, err := grpc.RegisterGrpc()
	if err != nil {
		panic(err)
	}
	//grpc服务注册到etcd
	r, err := grpc.RegisterEtcdServer(ctx)
	if err != nil {
		panic(err)
	}

	stop := func() {
		gc.Stop()
		r.Stop()
	}

	srv.Run(stop)
}
