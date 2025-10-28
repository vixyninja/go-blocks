package logx_test

import (
	"context"
	"testing"

	logx "github.com/vixyninja/go-blocks/pkg/logx"
)

func TestFactory_Basic(t *testing.T) {
	ctx := context.Background()

	types := []logx.LoggerType{logx.LoggerTypeStd, logx.LoggerTypeLogrus, logx.LoggerTypeZap, logx.LoggerTypeZerolog}
	for _, tp := range types {
		lg := logx.NewLogger(logx.LoggerConfig{Type: tp, Format: logx.FormatJSON})
		if lg == nil {
			t.Fatalf("failed to create logger: %s", tp)
		}
		lg.Info(ctx, "hello from %s", tp)
		lg.With(ctx, map[string]any{"k": "v"}).Info(ctx, "child")
	}
}
