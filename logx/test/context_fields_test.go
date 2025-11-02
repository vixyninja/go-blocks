package logx_test

import (
	"context"
	"testing"

	logx "github.com/vixyninja/go-blocks/logx"
)

func TestContextFields(t *testing.T) {
	ctx := context.Background()

	type Key string
	type Value any

	ctx = context.WithValue(ctx, Key("request_id"), Value("req-123"))
	ctx = context.WithValue(ctx, Key("user_id"), Value("user-456"))
	ctx = context.WithValue(ctx, Key("fields"), Value(map[string]any{"service": "test-service"}))

	logger := logx.NewZerologJSONLogger()
	logger.Info(ctx, "message with context fields")

	child := logger.With(ctx, map[string]any{"component": "auth"})
	child.Info(ctx, "child with context")
}
