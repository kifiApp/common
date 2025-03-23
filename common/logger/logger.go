package logger

import (
	"context"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"
	"log/slog"
	"time"
)

func CreateLogger(ctx context.Context, serviceName string, logLevel string) func() {
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

	tp := trace.NewTracerProvider(sdktrace.WithBatcher(exp,
		sdktrace.WithMaxExportBatchSize(512),
		sdktrace.WithBatchTimeout(5*time.Second)), sdktrace.WithResource(res))

	otel.SetTracerProvider(tp)
	return func() {
		if fErr := exp.Shutdown(ctx); fErr != nil {
			panic(fErr)
		}
	}
}
