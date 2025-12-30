package tracer

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc/stats"
)

func JaegerClientHandler(endpoints, name, env string, isOpenOnlySamplerError bool) (stats.Handler, error) {
	tp, tpErr := JaegerTraceProvider(
		endpoints,
		name,
		env,
		isOpenOnlySamplerError,
	)
	if tpErr != nil {
		return nil, tpErr
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	otelHandler := otelgrpc.NewClientHandler(
		otelgrpc.WithTracerProvider(tp), // 关联Jaeger追踪器
	)

	return otelHandler, nil
}

func JaegerServerHandler(endpoints, name, env string, isOpenOnlySamplerError bool) (stats.Handler, error) {
	tp, tpErr := JaegerTraceProvider(
		endpoints,
		name,
		env,
		isOpenOnlySamplerError,
	)
	if tpErr != nil {
		return nil, tpErr
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	otelHandler := otelgrpc.NewServerHandler(
		otelgrpc.WithTracerProvider(tp), // 关联Jaeger追踪器
	)

	return otelHandler, nil
}
