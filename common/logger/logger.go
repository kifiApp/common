package logger

import (
	"fmt"
	"log/slog"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
	"go.uber.org/zap/zapcore"
)

var once sync.Once
var logger *slog.Logger
var zapLogger *zap.Logger

func CreateLogger(serviceName string, logLevel string) (*slog.Logger, *zap.Logger) {
	once.Do(func() {
		zapLogger = configureZapLogger(serviceName, logLevel)
		logger = slog.New(zapslog.NewHandler(zapLogger.Core(), nil))
	})
	return logger, zapLogger
}

func Logger() *slog.Logger {
	once.Do(func() {
		logger = slog.New(zapslog.NewHandler(zap.L().Core(), nil))
	})
	return logger
}

func configureZapLogger(serviceName, logLevel string) *zap.Logger {
	var config zap.Config
	var debug bool
	if logLevel == "PRODUCTION" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		debug = true
	}

	config.EncoderConfig.TimeKey = "sts"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.DisableStacktrace = debug
	config.Development = debug
	config.InitialFields = map[string]interface{}{
		"service": serviceName,
	}
	if debug {
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	} else {
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}
	l, err := config.Build()
	if err != nil {
		panic(fmt.Sprintf("failed to build production logger: %v", err))
	}
	return l
}
