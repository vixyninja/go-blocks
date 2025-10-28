package logx_test

import (
	"context"
	"path/filepath"
	"testing"

	logx "github.com/vixyninja/go-blocks/pkg/logx"
)

func TestZerologLogger_Basic(t *testing.T) {
	ctx := context.Background()
	logger := logx.NewZerologJSONLogger()
	logger.Info(ctx, "info")
}

func TestZerologLogger_Rotation_JSON(t *testing.T) {
	ctx := context.Background()
	filePath := filepath.Join("logs", "zerolog_json.log")
	logger, err := logx.ZerologLoggerWithRotation(filePath, nil, logx.FormatJSON)
	if err != nil {
		t.Fatalf("failed to create zerolog rotation logger: %v", err)
	}
	logger.Info(ctx, "json info")
}

func TestZerologLogger_Rotation_Text(t *testing.T) {
	ctx := context.Background()
	filePath := filepath.Join("logs", "zerolog_text.log")
	logger, err := logx.ZerologLoggerWithRotation(filePath, nil, logx.FormatText)
	if err != nil {
		t.Fatalf("failed to create zerolog rotation logger: %v", err)
	}
	logger.Info(ctx, "text info")
}
