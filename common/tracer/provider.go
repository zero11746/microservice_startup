package tracer

import (
	"context"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.9.0"
	"go.opentelemetry.io/otel/trace"
)

func JaegerTraceProvider(endpoints, serviceName, environmentKey string, isOpenOnlyErrSampler bool) (*sdktrace.TracerProvider, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoints)))
	if err != nil {
		return nil, err
	}

	batchProcessor := sdktrace.NewBatchSpanProcessor(exp)
	var spanProcessor sdktrace.SpanProcessor

	if isOpenOnlyErrSampler {
		spanProcessor = NewErrorOnlySpanProcessor(batchProcessor)
	} else {
		spanProcessor = batchProcessor
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(spanProcessor),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			semconv.DeploymentEnvironmentKey.String(environmentKey),
		)),
	)

	return tp, nil
}

func GetTraceIDs(ctx context.Context) (traceID string, spanID string) {
	// 从 context 中获取当前 Span
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return "", ""
	}

	// 获取 Span 的上下文信息（包含 traceID 和 spanID）
	spanCtx := span.SpanContext()

	traceID = spanCtx.TraceID().String()
	spanID = spanCtx.SpanID().String()

	return traceID, spanID
}
