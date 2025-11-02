package logx

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

type ZerologLogger struct {
	logger zerolog.Logger
	fields map[string]any
}

type ZerologConfig struct {
	Output io.Writer
	Level  *zerolog.Level
}

func NewZerologLoggerWithLevel(level zerolog.Level) *ZerologLogger {
	config := ZerologConfig{
		Level: &level,
	}
	return NewZerologLoggerWithConfig(config)
}

func NewZerologLoggerWithOutput(output io.Writer) *ZerologLogger {
	config := ZerologConfig{
		Output: output,
	}
	return NewZerologLoggerWithConfig(config)
}

func NewZerologLogger() *ZerologLogger {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})

	return &ZerologLogger{
		logger: logger,
		fields: make(map[string]any),
	}
}

func NewZerologLoggerWithConfig(config ZerologConfig) *ZerologLogger {
	var logger zerolog.Logger

	if config.Output != nil {
		logger = zerolog.New(config.Output)
	} else {
		logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	}

	if config.Level != nil {
		logger = logger.Level(*config.Level)
	} else {
		logger = logger.Level(zerolog.DebugLevel)
	}

	if config.Level != nil {
		zerolog.SetGlobalLevel(*config.Level)
	}

	return &ZerologLogger{
		logger: logger,
		fields: make(map[string]any),
	}
}

func NewZerologJSONLogger() *ZerologLogger {
	config := ZerologConfig{
		Output: os.Stdout,
		Level:  &[]zerolog.Level{zerolog.DebugLevel}[0],
	}
	return NewZerologLoggerWithConfig(config)
}

func NewZerologConsoleLogger() *ZerologLogger {
	config := ZerologConfig{
		Output: zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339},
		Level:  &[]zerolog.Level{zerolog.DebugLevel}[0],
	}
	return NewZerologLoggerWithConfig(config)
}

func ZerologLoggerWithRotation(filePath string, rotConfig *RotationConfig, format Format) (*ZerologLogger, error) {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("[pkg.logx.ZerologLoggerWithRotation] failed to create log directory: %w", err)
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

	var out io.Writer = lumberjackLogger
	switch format {
	case FormatText:
		out = zerolog.ConsoleWriter{Out: lumberjackLogger, TimeFormat: time.RFC3339, NoColor: true}
	case FormatJSON:
		out = lumberjackLogger
	}

	logger := zerolog.New(out).With().Timestamp().Logger()
	logger = logger.Level(zerolog.DebugLevel)

	return &ZerologLogger{
		logger: logger,
		fields: make(map[string]any),
	}, nil
}

func (l *ZerologLogger) Debug(ctx context.Context, msg string, args ...any) {
	l.log(ctx, zerolog.DebugLevel, msg, args...)
}

func (l *ZerologLogger) Info(ctx context.Context, msg string, args ...any) {
	l.log(ctx, zerolog.InfoLevel, msg, args...)
}

func (l *ZerologLogger) Warn(ctx context.Context, msg string, args ...any) {
	l.log(ctx, zerolog.WarnLevel, msg, args...)
}

func (l *ZerologLogger) Error(ctx context.Context, msg string, args ...any) {
	l.log(ctx, zerolog.ErrorLevel, msg, args...)
}

func (l *ZerologLogger) Fatal(ctx context.Context, msg string, args ...any) {
	l.log(ctx, zerolog.FatalLevel, msg, args...)
	os.Exit(1)
}

func (l *ZerologLogger) With(_ context.Context, fields map[string]any) Logx {
	newFields := make(map[string]any)
	for k, v := range l.fields {
		newFields[k] = v
	}
	for k, v := range fields {
		newFields[k] = v
	}

	return &ZerologLogger{
		logger: l.logger,
		fields: newFields,
	}
}

func (l *ZerologLogger) log(ctx context.Context, level zerolog.Level, msg string, args ...any) {
	event := l.logger.WithLevel(level)

	for k, v := range l.fields {
		event = event.Interface(k, v)
	}

	contextFields := extractContextFields(ctx)
	for k, v := range contextFields {
		event = event.Interface(k, v)
	}

	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	event.Msg(msg)
}
