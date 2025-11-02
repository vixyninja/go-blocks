package logx_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	logx "github.com/vixyninja/go-blocks/logx"
)

func TestRotation_AllBackends(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		name     string
		filename string
		create   func(string) (logx.Logx, error)
	}{
		{"std_text", "std_rotation.log", func(p string) (logx.Logx, error) { return logx.StdLoggerWithRotation(p, nil) }},
		{"logrus_json", "logrus_rotation.log", func(p string) (logx.Logx, error) { return logx.LogrusLoggerWithRotation(p, nil, logx.FormatJSON) }},
		{"logrus_text", "logrus_rotation_text.log", func(p string) (logx.Logx, error) { return logx.LogrusLoggerWithRotation(p, nil, logx.FormatText) }},
		{"zap_json", "zap_rotation.log", func(p string) (logx.Logx, error) { return logx.ZapLoggerWithRotation(p, nil, logx.FormatJSON) }},
		{"zap_text", "zap_rotation_text.log", func(p string) (logx.Logx, error) { return logx.ZapLoggerWithRotation(p, nil, logx.FormatText) }},
		{"zerolog_json", "zerolog_rotation.log", func(p string) (logx.Logx, error) { return logx.ZerologLoggerWithRotation(p, nil, logx.FormatJSON) }},
		{"zerolog_text", "zerolog_rotation_text.log", func(p string) (logx.Logx, error) { return logx.ZerologLoggerWithRotation(p, nil, logx.FormatText) }},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			path := filepath.Join("logs", cs.filename)
			logger, err := cs.create(path)
			if err != nil {
				t.Fatalf("create logger: %v", err)
			}
			for i := 0; i < 10; i++ {
				logger.Info(ctx, "rotation test %d", i)
			}
			if info, err := os.Stat(path); err != nil || info.Size() == 0 {
				t.Fatalf("log file not created or empty: %v", err)
			}
		})
	}
}
