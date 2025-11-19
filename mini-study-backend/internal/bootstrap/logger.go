package bootstrap

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLogger builds a zap.Logger based on config log level.
func InitLogger(cfg *Config) (*zap.Logger, error) {
	level := zapcore.InfoLevel
	if err := level.Set(strings.ToLower(cfg.App.LogLevel)); err != nil {
		level = zapcore.InfoLevel
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	zapCfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Development:      cfg.App.Env == "dev" || cfg.App.Env == "local",
		Encoding:         "json",
		EncoderConfig:    encoderCfg,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	if cfg.App.Env == "dev" || cfg.App.Env == "local" {
		zapCfg.Encoding = "console"
	}

	return zapCfg.Build()
}
