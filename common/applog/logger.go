package applog

import (
	"context"
	"fmt"
)

type ctxKey string

const (
	Trace      ctxKey = "trace"
	SpanID     ctxKey = "spanID"
	TraceIDKey ctxKey = "traceID"
)

var logger *Logger

func GetLoggerInstance() *Logger {
	return logger
}

// InitLoggers 初始化日志
func InitLoggers(log LogConfig) error {
	client := NewLogger(&log)
	if client == nil {
		return fmt.Errorf("logger is nil")
	}
	logger = client
	return nil
}

// WrapGDPLogger 获取日志
func WrapGDPLogger(ctx context.Context) *Tracer {
	// 从context中获取tracer
	tracer := ctx.Value(Trace)
	tracerClient, ok := tracer.(*Tracer)

	if tracer == nil || !ok {
		tracerClient = NewTracer(logger, "", "")
	}

	return tracerClient
}
