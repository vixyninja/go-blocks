package logx

type LoggerType string
type Format string

type LoggerConfig struct {
	Type   LoggerType
	Format Format // "json" or "text"
	Output string // "stdout", "stderr", or file path
	Prefix string // for std logger
}

const (
	LoggerTypeStd     LoggerType = "std"
	LoggerTypeLogrus  LoggerType = "logrus"
	LoggerTypeZap     LoggerType = "zap"
	LoggerTypeZerolog LoggerType = "zerolog"
)

const (
	FormatJSON Format = "json"
	FormatText Format = "text"
)

const (
	MaxSize    int  = 100
	MaxAge     int  = 30
	MaxBackups int  = 3
	Compress   bool = true
	LocalTime  bool = false
)

const (
	TimestampFormat = "2006-01-02T15:04:05.000Z07:00"
	TimeKey         = "timestamp"
)

type RotationConfig struct {
	MaxSize    int  // Maximum size in megabytes before rotating (default: 100MB)
	MaxAge     int  // Maximum number of days to retain old log files (default: 0 = unlimited)
	MaxBackups int  // Maximum number of old log files to retain (default: 0 = unlimited)
	Compress   bool // Compress rotated log files (default: false)
	LocalTime  bool // Use local time for rotation (default: false)
}

func DefaultRotationConfig() *RotationConfig {
	return &RotationConfig{
		MaxSize:    MaxSize,
		MaxAge:     MaxAge,
		MaxBackups: MaxBackups,
		Compress:   Compress,
		LocalTime:  LocalTime,
	}
}

func NewLogger(config LoggerConfig) Logx {
	switch config.Type {
	case LoggerTypeStd:
		if config.Prefix != "" {
			return NewStdLoggerWithPrefix(config.Prefix)
		}
		return NewStdLogger()

	case LoggerTypeLogrus:
		if config.Format == FormatJSON {
			return NewLogrusJSONLogger()
		}
		return NewLogrusTextLogger()

	case LoggerTypeZap:
		if config.Format == FormatJSON {
			return NewZapJSONLogger()
		}
		return NewZapConsoleLogger()

	case LoggerTypeZerolog:
		if config.Format == FormatJSON {
			return NewZerologJSONLogger()
		}
		return NewZerologConsoleLogger()

	default:
		return NewStdLogger()
	}
}

func NewDefaultLogger() Logx {
	return NewLogger(LoggerConfig{
		Type:   LoggerTypeStd,
		Format: FormatText,
		Output: "stdout",
	})
}

func NewJSONLogger() Logx {
	return NewLogger(LoggerConfig{
		Type:   LoggerTypeZerolog,
		Format: FormatJSON,
		Output: "stdout",
	})
}

func NewTextLogger() Logx {
	return NewLogger(LoggerConfig{
		Type:   LoggerTypeLogrus,
		Format: FormatText,
		Output: "stdout",
	})
}

func NewHighPerformanceLogger() Logx {
	return NewLogger(LoggerConfig{
		Type:   LoggerTypeZap,
		Format: FormatJSON,
		Output: "stdout",
	})
}

func NewSimpleLogger() Logx {
	return NewLogger(LoggerConfig{
		Type:   LoggerTypeStd,
		Format: FormatText,
		Output: "stdout",
	})
}
