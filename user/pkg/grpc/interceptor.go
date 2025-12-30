package grpc

import (
	"common/applog"
	"common/errs"
	"common/tracer"
	"context"
	"google.golang.org/grpc"
	"user/pkg/errors"
)

func TraceIDInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		tracerID, spanID := tracer.GetTraceIDs(ctx)
		t := applog.NewTracer(applog.GetLoggerInstance(), tracerID, spanID)
		newCtx := context.WithValue(ctx, applog.Trace, t)
		return handler(newCtx, req)
	}
}

func ErrorLogInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// 执行 service 方法（先调用 handler，获取结果和错误）
		resp, err = handler(ctx, req)

		logger := applog.WrapGDPLogger(ctx)

		// 若有错误，记录日志
		if err != nil {
			_, errorMsg := errs.ParseGrpcError(err)
			logger.WithReq(req).WithResp(resp).Error(errorMsg)
		}

		return resp, err
	}
}

func ErrorInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, err = handler(ctx, req)

		code, _ := errs.ParseGrpcError(err)
		if _, ok := errors.Errors[errs.ErrorCode(code)]; ok {
			return resp, errors.Errors[errs.ErrorCode(code)]
		} else {
			return resp, errors.UnknownError
		}
	}
}
