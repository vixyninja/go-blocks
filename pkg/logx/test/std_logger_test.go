package logx_test

import (
	"context"
	"path/filepath"
	"testing"

	logx "github.com/vixyninja/go-blocks/pkg/logx"
)

func TestStdLogger_Basic(t *testing.T) {
	ctx := context.Background()
	logger := logx.NewStdLogger()

	logger.Debug(ctx, "debug message")
	logger.Info(ctx, "info message")
	logger.Warn(ctx, "warn message")
	logger.Error(ctx, "error message")

	child := logger.With(ctx, map[string]any{"component": "test"})
	child.Info(ctx, "child message")
}

func TestStdLogger_Rotation(t *testing.T) {
	ctx := context.Background()
	filePath := filepath.Join("logs", "std_text.log")
	logger, err := logx.StdLoggerWithRotation(filePath, nil)
	if err != nil {
		t.Fatalf("failed to create std rotation logger: %v", err)
	}

	logger.Info(ctx, "info to file")
	logger.Warn(ctx, "warn to file")
}
