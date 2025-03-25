package logger

import (
	"context"
	"fmt"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"
	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
	"go.uber.org/zap/zapcore"
	"log/slog"
	"sync"
	"time"
)

var once sync.Once
var logger *slog.Logger
var zapLogger *zap.Logger

func CreateLogger(serviceName string, logLevel string) (*slog.Logger, *zap.Logger) {
	once.Do(func() {
		zapLogger = configureZapLogger(serviceName, logLevel)
		logger = slog.New(zapslog.NewHandler(zapLogger.Core()))
	})
	return logger, zapLogger
}

func CreateLoggerWithMetrics(ctx context.Context, serviceName string, logLevel string) func() {
	logExporter, err := otlploggrpc.New(ctx, otlploggrpc.WithEndpoint(""), otlploggrpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	res, err := resource.New(ctx, resource.WithAttributes(semconv.ServiceNameKey.String(serviceName)))
	if err != nil {
		panic(err)
	}

	lp := log.NewLoggerProvider(log.WithResource(res), log.WithProcessor(log.NewBatchProcessor(logExporter)))

	global.SetLoggerProvider(lp)
	slog.SetDefault(otelslog.NewLogger(serviceName))
	return func() {
		if fErr := logExporter.Shutdown(ctx); fErr != nil {
			panic(fErr)
		}
	}
}

func CreateTracer(ctx context.Context, serviceName string) func() {
	exp, err := otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint(""), otlptracegrpc.WithInsecure())

	if err != nil {
		panic(err)
	}

	res, err := resource.New(ctx, resource.WithAttributes(semconv.ServiceNameKey.String(serviceName)))

	if err != nil {
		panic(err)
	}

	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exp,
		sdktrace.WithMaxExportBatchSize(512),
		sdktrace.WithBatchTimeout(5*time.Second)), sdktrace.WithResource(res))

	otel.SetTracerProvider(tp)
	return func() {
		if fErr := exp.Shutdown(ctx); fErr != nil {
			panic(fErr)
		}
	}
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
