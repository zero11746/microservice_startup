package grpc

import (
	"common/applog"
	discovery2 "common/discovery"
	"common/tracer"
	"context"
	"errors"
	"fmt"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	userservice "grpc/user/user"
	"log"
	"net"
	"user/config"
	"user/internal/service"
)

type gRPCConfig struct {
	Addr         string
	RegisterFunc func(*grpc.Server)
}

func RegisterGrpc() (*grpc.Server, error) {
	otelHandler, err := tracer.JaegerServerHandler(
		config.GetConfig().Jaeger.Endpoints,
		config.GetConfig().Server.Name,
		config.GetConfig().Server.Env,
		config.GetConfig().Jaeger.IsOpenOnlySamplerError,
	)
	if err != nil {
		return nil, err
	}

	// 创建gRPC服务器
	s := grpc.NewServer(
		//grpc.Creds(), // 使用TLS
		grpc.StatsHandler(otelHandler),
		grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
			// 注册其他拦截器
			TraceIDInterceptor(),
			ErrorLogInterceptor(),
			ErrorInterceptor(),
		)),
	)

	// 注册gRPC业务服务
	c := gRPCConfig{
		Addr: config.GetConfig().Grpc.Addr,
		RegisterFunc: func(s *grpc.Server) {
			// 注册服务
			userservice.RegisterUserServer(s, service.NewUserService())
		},
	}
	c.RegisterFunc(s)

	// 启动gRPC服务监听
	lis, err := net.Listen("tcp", c.Addr)
	if err != nil {
		return nil, fmt.Errorf("listen tcp %s error: %w", c.Addr, err)
	}

	go func() {
		if err := s.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		}
	}()

	return s, nil
}

func RegisterEtcdServer(ctx context.Context) (*discovery2.Register, error) {
	info := discovery2.Server{
		Name:    config.GetConfig().Grpc.Name,
		Addr:    config.GetConfig().Grpc.EtcdAddr,
		Version: config.GetConfig().Grpc.Version,
		Weight:  config.GetConfig().Grpc.Weight,
	}

	r := discovery2.NewRegister(config.GetConfig().Etcd.Addrs, applog.WrapGDPLogger(ctx))
	r, _, err := r.Register(info, 2)
	if err != nil {
		log.Fatalln(err)
	}

	return r, nil
}
