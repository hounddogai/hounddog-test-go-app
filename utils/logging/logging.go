package logging

import (
	"context"
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type spanKey struct{}

type Span struct {
	TraceID string
	SpanID  string
	Baggage map[string]string
}

type CustomLogger struct{}

var MyCustomlog = &CustomLogger{}

var (
	loggerOnce sync.Once
	baseLogger *zap.Logger
)

func (c *CustomLogger) Infof(ctx context.Context, msg string, args ...any) {
	globalLogCtx(ctx).Sugar().Infof(msg, args...)
}

func ContextWithSpan(ctx context.Context, span Span) context.Context {
	return context.WithValue(ctx, spanKey{}, span)
}

func spanFromCtx(ctx context.Context) (Span, bool) {
	span, ok := ctx.Value(spanKey{}).(Span)
	return span, ok
}

func globalLogCtx(ctx context.Context) *zap.Logger {
	span, _ := spanFromCtx(ctx)
	return withSpan(span, globalLog())
}

func withSpan(span Span, logger *zap.Logger) *zap.Logger {
	fields := []zap.Field{}
	if span.TraceID != "" {
		fields = append(fields, zap.String("trace_id", span.TraceID))
	}
	if span.SpanID != "" {
		fields = append(fields, zap.String("span_id", span.SpanID))
	}
	if len(span.Baggage) > 0 {
		for k, v := range span.Baggage {
			fields = append(fields, zap.String("baggage_"+k, v))
		}
	}
	return logger.With(fields...)
}

func globalLog() *zap.Logger {
	loggerOnce.Do(func() {
		encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
		core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.InfoLevel)
		baseLogger = zap.New(&prefixedCore{core: core})
	})
	return baseLogger
}

// prefixedCore sanitizes message fields before forwarding to the underlying core.
type prefixedCore struct {
	core zapcore.Core
}

func (p *prefixedCore) Enabled(lvl zapcore.Level) bool { return p.core.Enabled(lvl) }

func (p *prefixedCore) With(fields []zapcore.Field) zapcore.Core {
	return &prefixedCore{core: p.core.With(fields)}
}

func (p *prefixedCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if p.Enabled(ent.Level) {
		return p.core.Check(ent, ce)
	}
	return ce
}

func (p *prefixedCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	ent.Message = sanitize(ent.Message)
	return p.core.Write(ent, fields)
}

func (p *prefixedCore) Sync() error { return p.core.Sync() }

func sanitize(msg string) string {
	return strings.ReplaceAll(msg, "secret", "[redacted]")
}
