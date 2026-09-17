/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

// Package tracing wires OpenTelemetry tracing into the shared server and
// client constructors. It is inert unless OTEL_EXPORTER_OTLP_ENDPOINT is set,
// which the OpenTelemetry Operator injects into annotated namespaces; local
// runs and unit tests keep the global no-op tracer.
package tracing

import (
	"context"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// EndpointEnv is the standard OTLP endpoint variable; its presence switches
// tracing on.
const EndpointEnv = "OTEL_EXPORTER_OTLP_ENDPOINT"

var (
	once     sync.Once
	provider *sdktrace.TracerProvider
)

// Init installs the global tracer provider and W3C propagators when
// OTEL_EXPORTER_OTLP_ENDPOINT is set. Service name, resource attributes and
// sampler come from the standard OTEL_* environment variables. Only the first
// call does anything, so every constructor may call it.
func Init() {
	once.Do(func() {
		endpoint := os.Getenv(EndpointEnv)
		if endpoint == "" {
			return
		}

		ctx := context.Background()

		exporter, err := otlptracegrpc.New(ctx)
		if err != nil {
			log.Errorf("Tracing disabled: OTLP exporter init failed: %v", err)

			return
		}

		res, err := resource.New(ctx,
			resource.WithFromEnv(),
			resource.WithTelemetrySDK(),
			resource.WithHost(),
		)
		if err != nil {
			log.Warnf("Tracing resource detection incomplete: %v", err)
		}

		provider = sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(res),
		)

		otel.SetTracerProvider(provider)
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		))

		log.AddHook(LogHook{})

		log.Infof("Tracing enabled, exporting to %s", endpoint)
	})
}

// Shutdown flushes buffered spans. Services that handle SIGTERM can call it
// before exiting; it is a no-op when tracing is off.
func Shutdown(ctx context.Context) error {
	if provider == nil {
		return nil
	}

	return provider.Shutdown(ctx)
}

// LogHook adds trace_id and span_id to logrus entries that carry a context
// with an active span (log.WithContext(ctx)).
type LogHook struct{}

func (LogHook) Levels() []log.Level {
	return log.AllLevels
}

func (LogHook) Fire(entry *log.Entry) error {
	if entry.Context == nil {
		return nil
	}

	sc := trace.SpanContextFromContext(entry.Context)
	if !sc.IsValid() {
		return nil
	}

	entry.Data["trace_id"] = sc.TraceID().String()
	entry.Data["span_id"] = sc.SpanID().String()

	return nil
}
