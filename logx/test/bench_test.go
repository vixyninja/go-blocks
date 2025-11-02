package logx_test

import (
	"context"
	"testing"

	logx "github.com/vixyninja/go-blocks/logx"
)

func BenchmarkStdLogger(b *testing.B) {
	ctx := context.Background()
	logger := logx.NewStdLogger()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info(ctx, "bench %d", i)
	}
}

func BenchmarkLogrusLogger(b *testing.B) {
	ctx := context.Background()
	logger := logx.NewLogrusJSONLogger()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info(ctx, "bench %d", i)
	}
}

func BenchmarkZapLogger(b *testing.B) {
	ctx := context.Background()
	logger := logx.NewZapJSONLogger()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info(ctx, "bench %d", i)
	}
}

func BenchmarkZerologLogger(b *testing.B) {
	ctx := context.Background()
	logger := logx.NewZerologJSONLogger()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info(ctx, "bench %d", i)
	}
}
