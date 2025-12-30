package tracer

import (
	"context"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type ErrorOnlySpanProcessor struct {
	next sdktrace.SpanProcessor
}

func NewErrorOnlySpanProcessor(next sdktrace.SpanProcessor) *ErrorOnlySpanProcessor {
	return &ErrorOnlySpanProcessor{next: next}
}

// OnStart 开始 Span 时，委托给下一个处理器
func (p *ErrorOnlySpanProcessor) OnStart(ctx context.Context, s sdktrace.ReadWriteSpan) {
	p.next.OnStart(ctx, s)
}

// OnEnd 结束 Span 时，判断错误并导出
func (p *ErrorOnlySpanProcessor) OnEnd(s sdktrace.ReadOnlySpan) {
	hasError := false

	for _, kv := range s.Attributes() {
		if kv.Key == "error" && kv.Value.AsBool() {
			hasError = true
			break
		}
	}

	// 检查状态码
	if !hasError {
		if s.Status().Code == codes.Error {
			hasError = true
		}
	}

	if hasError {
		p.next.OnEnd(s)
	}
}

// Shutdown 关闭处理器
func (p *ErrorOnlySpanProcessor) Shutdown(ctx context.Context) error {
	return p.next.Shutdown(ctx)
}

// ForceFlush 强制刷新
func (p *ErrorOnlySpanProcessor) ForceFlush(ctx context.Context) error {
	return p.next.ForceFlush(ctx)
}
