package logx_test

import (
	"context"
	"path/filepath"
	"testing"

	logx "github.com/vixyninja/go-blocks/pkg/logx"
)

func TestLogrusLogger_Basic(t *testing.T) {
	ctx := context.Background()
	logger := logx.NewLogrusJSONLogger()
	logger.Info(ctx, "info")
	logger.With(ctx, map[string]any{"k": "v"}).Warn(ctx, "warn")
}

func TestLogrusLogger_Rotation_Text(t *testing.T) {
	ctx := context.Background()
	filePath := filepath.Join("logs", "logrus_text.log")
	logger, err := logx.LogrusLoggerWithRotation(filePath, nil, logx.FormatText)
	if err != nil {
		t.Fatalf("failed to create logrus rotation logger: %v", err)
	}
	logger.Info(ctx, "text info")
}

func TestLogrusLogger_Rotation_JSON(t *testing.T) {
	ctx := context.Background()
	filePath := filepath.Join("logs", "logrus_json.log")
	logger, err := logx.LogrusLoggerWithRotation(filePath, nil, logx.FormatJSON)
	if err != nil {
		t.Fatalf("failed to create logrus rotation logger: %v", err)
	}
	logger.Info(ctx, "json info")
}
