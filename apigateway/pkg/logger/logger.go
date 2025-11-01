package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(env string) *zap.Logger {
	var cfg zap.Config

	if env == "production" {
		cfg = zap.Config{
			Level:       zap.NewAtomicLevelAt(zapcore.InfoLevel),
			Encoding:    "json",
			OutputPaths: []string{"stdout"},
			EncoderConfig: zapcore.EncoderConfig{
				TimeKey:        "timestamp",
				LevelKey:       "level",
				NameKey:        "logger",
				CallerKey:      "caller",
				MessageKey:     "message",
				StacktraceKey:  "stacktrace",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.CapitalLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.StringDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			},
		}
	} else {
		cfg = zap.Config{
			Level:       zap.NewAtomicLevelAt(zapcore.DebugLevel),
			Development: true,
			Encoding:    "console",
			OutputPaths: []string{"stdout"},
			EncoderConfig: zapcore.EncoderConfig{
				TimeKey:        "T",
				LevelKey:       "L",
				CallerKey:      "C",
				MessageKey:     "M",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.CapitalColorLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.StringDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			},
		}
	}

	logger, err := cfg.Build()
	if err != nil {
		fallback := zap.NewExample()
		fallback.Error("Failed to build zap logger", zap.Error(err))
		return fallback
	}

	zap.ReplaceGlobals(logger)
	return logger
}

func Sync(logger *zap.Logger) {
	_ = logger.Sync()
}

func InitGlobal(env string) *zap.Logger {
	logger := New(env)
	zap.ReplaceGlobals(logger)
	return logger
}
