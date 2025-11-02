package logx_test

import (
	"context"
	"path/filepath"
	"testing"

	logx "github.com/vixyninja/go-blocks/logx"
)

func TestZapLogger_Basic(t *testing.T) {
	ctx := context.Background()
	logger := logx.NewZapJSONLogger()
	logger.Info(ctx, "info")
	_ = logger.Sync()
}

func TestZapLogger_Rotation_JSON(t *testing.T) {
	ctx := context.Background()
	filePath := filepath.Join("logs", "zap_json.log")
	logger, err := logx.ZapLoggerWithRotation(filePath, nil, logx.FormatJSON)
	if err != nil {
		t.Fatalf("failed to create zap rotation logger: %v", err)
	}
	logger.Info(ctx, "json info")
}

func TestZapLogger_Rotation_Text(t *testing.T) {
	ctx := context.Background()
	filePath := filepath.Join("logs", "zap_text.log")
	logger, err := logx.ZapLoggerWithRotation(filePath, nil, logx.FormatText)
	if err != nil {
		t.Fatalf("failed to create zap rotation logger: %v", err)
	}
	logger.Info(ctx, "text info")
}
