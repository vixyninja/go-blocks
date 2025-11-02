package logx

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type ZapLogger struct {
	logger *zap.Logger
	fields map[string]any
}

func NewZapLogger() *ZapLogger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.TimeKey = TimeKey
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, _ := config.Build()

	return &ZapLogger{
		logger: logger,
		fields: make(map[string]any),
	}
}

func NewZapLoggerWithConfig(config ZapConfig) *ZapLogger {
	var logger *zap.Logger
	var err error

	if config.Production {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		logger, _ = zap.NewDevelopment()
	}

	if config.Level != nil {
		atomicLevel := zap.NewAtomicLevelAt(*config.Level)
		logger = logger.WithOptions(zap.IncreaseLevel(atomicLevel))
	}

	return &ZapLogger{
		logger: logger,
		fields: make(map[string]any),
	}
}

func NewZapJSONLogger() *ZapLogger {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = TimeKey
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, _ := config.Build()

	return &ZapLogger{
		logger: logger,
		fields: make(map[string]any),
	}
}

func NewZapConsoleLogger() *ZapLogger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.TimeKey = TimeKey
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger, _ := config.Build()

	return &ZapLogger{
		logger: logger,
		fields: make(map[string]any),
	}
}

func (l *ZapLogger) Debug(ctx context.Context, msg string, args ...any) {
	l.log(ctx, zap.DebugLevel, msg, args...)
}

func (l *ZapLogger) Info(ctx context.Context, msg string, args ...any) {
	l.log(ctx, zap.InfoLevel, msg, args...)
}

func (l *ZapLogger) Warn(ctx context.Context, msg string, args ...any) {
	l.log(ctx, zap.WarnLevel, msg, args...)
}

func (l *ZapLogger) Error(ctx context.Context, msg string, args ...any) {
	l.log(ctx, zap.ErrorLevel, msg, args...)
}

func (l *ZapLogger) Fatal(ctx context.Context, msg string, args ...any) {
	l.log(ctx, zap.FatalLevel, msg, args...)
	os.Exit(1)
}

func (l *ZapLogger) With(_ context.Context, fields map[string]any) Logx {
	newFields := make(map[string]any)
	for k, v := range l.fields {
		newFields[k] = v
	}
	for k, v := range fields {
		newFields[k] = v
	}

	return &ZapLogger{
		logger: l.logger,
		fields: newFields,
	}
}

func ZapLoggerWithRotation(filePath string, rotConfig *RotationConfig, format Format) (*ZapLogger, error) {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("[pkg.logx.ZapLoggerWithRotation] failed to create log directory: %w", err)
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

	var encoder zapcore.Encoder
	switch format {
	case FormatJSON:
		encCfg := zap.NewProductionEncoderConfig()
		encCfg.TimeKey = TimeKey
		encCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder = zapcore.NewJSONEncoder(encCfg)
	case FormatText:
		encCfg := zap.NewDevelopmentEncoderConfig()
		encCfg.TimeKey = TimeKey
		encCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder = zapcore.NewConsoleEncoder(encCfg)
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(lumberjackLogger),
		zap.NewAtomicLevelAt(zap.DebugLevel),
	)

	return &ZapLogger{
		logger: zap.New(core),
		fields: make(map[string]any),
	}, nil
}

func (l *ZapLogger) log(ctx context.Context, level zapcore.Level, msg string, args ...any) {
	zapFields := l.convertToZapFields(l.fields)

	contextFields := extractContextFields(ctx)
	if len(contextFields) > 0 {
		contextZapFields := l.convertToZapFields(contextFields)
		zapFields = append(zapFields, contextZapFields...)
	}

	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	switch level {
	case zap.DebugLevel:
		l.logger.Debug(msg, zapFields...)
	case zap.InfoLevel:
		l.logger.Info(msg, zapFields...)
	case zap.WarnLevel:
		l.logger.Warn(msg, zapFields...)
	case zap.ErrorLevel:
		l.logger.Error(msg, zapFields...)
	case zap.FatalLevel:
		l.logger.Fatal(msg, zapFields...)
	case zap.DPanicLevel, zap.PanicLevel, zapcore.InvalidLevel:
		l.logger.Panic(msg, zapFields...)
	default:
		l.logger.Debug(msg, zapFields...)
	}
}

func (l *ZapLogger) convertToZapFields(fields map[string]any) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return zapFields
}

type ZapConfig struct {
	Production bool
	Level      *zapcore.Level
}

func (l *ZapLogger) Sync() error {
	return l.logger.Sync()
}

func (l *ZapLogger) Close() error {
	return l.logger.Sync()
}
