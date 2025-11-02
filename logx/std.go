package logx

import (
	"context"
	"fmt"
	"log"
	"maps"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

type StdLogger struct {
	logger *log.Logger
	fields map[string]any
}

func NewStdLogger() *StdLogger {
	return &StdLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile),
		fields: make(map[string]any),
	}
}

func NewStdLoggerWithPrefix(prefix string) *StdLogger {
	return &StdLogger{
		logger: log.New(os.Stdout, prefix, log.LstdFlags|log.Lshortfile),
		fields: make(map[string]any),
	}
}

func (l *StdLogger) Debug(ctx context.Context, msg string, args ...any) {
	l.log(ctx, LevelDebug, msg, args...)
}

func (l *StdLogger) Info(ctx context.Context, msg string, args ...any) {
	l.log(ctx, LevelInfo, msg, args...)
}

func (l *StdLogger) Warn(ctx context.Context, msg string, args ...any) {
	l.log(ctx, LevelWarn, msg, args...)
}

func (l *StdLogger) Error(ctx context.Context, msg string, args ...any) {
	l.log(ctx, LevelError, msg, args...)
}

func (l *StdLogger) Fatal(ctx context.Context, msg string, args ...any) {
	l.log(ctx, LevelFatal, msg, args...)
	os.Exit(1)
}

func (l *StdLogger) With(_ context.Context, fields map[string]any) Logx {
	newFields := make(map[string]any)
	for k, v := range l.fields {
		newFields[k] = v
	}
	for k, v := range fields {
		newFields[k] = v
	}

	return &StdLogger{
		logger: l.logger,
		fields: newFields,
	}
}

func StdLoggerWithRotation(filePath string, rotConfig *RotationConfig) (*StdLogger, error) {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("[pkg.logx.StdLoggerWithRotation] failed to create log directory: %w", err)
	}

	if rotConfig == nil {
		rotConfig = DefaultRotationConfig()
	}

	lumberjackLogger := &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    rotConfig.MaxSize,
		MaxBackups: rotConfig.MaxBackups,
		MaxAge:     rotConfig.MaxAge,
		Compress:   rotConfig.Compress,
		LocalTime:  rotConfig.LocalTime,
	}

	return &StdLogger{
		logger: log.New(lumberjackLogger, "", log.LstdFlags|log.Lshortfile),
		fields: make(map[string]any),
	}, nil
}

func (l *StdLogger) log(ctx context.Context, level LogLevel, msg string, args ...any) {
	formattedMsg := msg
	if len(args) > 0 {
		formattedMsg = fmt.Sprintf(msg, args...)
	}

	contextFields := extractContextFields(ctx)

	allFields := make(map[string]any)
	for k, v := range l.fields {
		allFields[k] = v
	}
	for k, v := range contextFields {
		allFields[k] = v
	}

	finalMsg := fmt.Sprintf("[%s] %s", level.String(), formattedMsg)
	if len(allFields) > 0 {
		finalMsg += " " + formatFields(allFields)
	}

	l.logger.Println(finalMsg)
}

func (l LogLevel) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

func extractContextFields(ctx context.Context) map[string]any {
	fields := make(map[string]any)

	if v := ctx.Value("fields"); v != nil {
		if f, ok := v.(map[string]any); ok {
			maps.Copy(fields, f)
		}
	}

	if reqID := ctx.Value("request_id"); reqID != nil {
		fields["request_id"] = reqID
	}

	if traceID := ctx.Value("trace_id"); traceID != nil {
		fields["trace_id"] = traceID
	}

	return fields
}

func formatFields(fields map[string]any) string {
	if len(fields) == 0 {
		return ""
	}

	var parts []string
	for k, v := range fields {
		parts = append(parts, fmt.Sprintf("%s=%v", k, v))
	}

	result := ""
	for i, part := range parts {
		if i > 0 {
			result += " "
		}
		result += part
	}

	return result
}
