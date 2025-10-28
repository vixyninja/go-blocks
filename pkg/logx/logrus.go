package logx

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogrusLogger struct {
	logger *logrus.Logger
	fields map[string]any
}

type LogrusConfig struct {
	Output    io.Writer
	Formatter logrus.Formatter
	Level     *logrus.Level
}

func NewLogrusJSONLogger() *LogrusLogger {
	config := LogrusConfig{
		Formatter: &logrus.JSONFormatter{
			TimestampFormat: TimestampFormat,
		},
	}
	return NewLogrusLoggerWithConfig(config)
}

func NewLogrusTextLogger() *LogrusLogger {
	config := LogrusConfig{
		Formatter: &logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: TimestampFormat,
		},
	}
	return NewLogrusLoggerWithConfig(config)
}

func NewLogrusLogger() *LogrusLogger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logger.SetLevel(logrus.DebugLevel)

	return &LogrusLogger{
		logger: logger,
		fields: make(map[string]any),
	}
}

func NewLogrusLoggerWithConfig(config LogrusConfig) *LogrusLogger {
	logger := logrus.New()

	if config.Output != nil {
		logger.SetOutput(config.Output)
	} else {
		logger.SetOutput(os.Stdout)
	}

	if config.Formatter != nil {
		logger.SetFormatter(config.Formatter)
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	if config.Level != nil {
		logger.SetLevel(*config.Level)
	} else {
		logger.SetLevel(logrus.DebugLevel)
	}

	return &LogrusLogger{
		logger: logger,
		fields: make(map[string]any),
	}
}

func (l *LogrusLogger) Debug(ctx context.Context, msg string, args ...any) {
	l.log(ctx, logrus.DebugLevel, msg, args...)
}

func (l *LogrusLogger) Info(ctx context.Context, msg string, args ...any) {
	l.log(ctx, logrus.InfoLevel, msg, args...)
}

func (l *LogrusLogger) Warn(ctx context.Context, msg string, args ...any) {
	l.log(ctx, logrus.WarnLevel, msg, args...)
}

func (l *LogrusLogger) Error(ctx context.Context, msg string, args ...any) {
	l.log(ctx, logrus.ErrorLevel, msg, args...)
}

func (l *LogrusLogger) Fatal(ctx context.Context, msg string, args ...any) {
	l.log(ctx, logrus.FatalLevel, msg, args...)
	os.Exit(1)
}

func (l *LogrusLogger) With(_ context.Context, fields map[string]any) Logx {
	newFields := make(map[string]any)
	for k, v := range l.fields {
		newFields[k] = v
	}
	for k, v := range fields {
		newFields[k] = v
	}

	return &LogrusLogger{
		logger: l.logger,
		fields: newFields,
	}
}

func (l *LogrusLogger) log(ctx context.Context, level logrus.Level, msg string, args ...any) {
	entry := l.logger.WithFields(logrus.Fields(l.fields))

	contextFields := extractContextFields(ctx)
	if len(contextFields) > 0 {
		entry = entry.WithFields(logrus.Fields(contextFields))
	}

	if len(args) > 0 {
		entry.Logf(level, msg, args...)
	} else {
		entry.Log(level, msg)
	}
}

func LogrusLoggerWithRotation(filePath string, rotConfig *RotationConfig, format Format) (*LogrusLogger, error) {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("[pkg.logx.LogrusLoggerWithRotation] failed to create log directory: %w", err)
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

	var formatter logrus.Formatter
	switch format {
	case FormatJSON:
		formatter = &logrus.JSONFormatter{TimestampFormat: TimestampFormat}
	case FormatText:
		formatter = &logrus.TextFormatter{FullTimestamp: true, TimestampFormat: TimestampFormat}
	}

	config := LogrusConfig{
		Output:    lumberjackLogger,
		Formatter: formatter,
	}
	return NewLogrusLoggerWithConfig(config), nil
}
