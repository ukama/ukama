/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// CheckServiceHealth dials addr and performs a health check against the
// standard grpc.health.v1.Health service.
//
// It returns nil when the service reports SERVING, and a descriptive error
// otherwise (unreachable, timed out, or reporting a non-serving status).
// The caller controls the overall deadline through ctx.
func CheckServiceHealth(ctx context.Context, addr string) error {
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to create connection to %s: %w", addr, err)
	}
	defer func() {
		_ = conn.Close()
	}()

	resp, err := healthpb.NewHealthClient(conn).Check(ctx,
		&healthpb.HealthCheckRequest{})
	if err == nil {
		if resp.GetStatus() != healthpb.HealthCheckResponse_SERVING {
			return fmt.Errorf("service at %s reported status %s",
				addr, resp.GetStatus().String())
		}

		return nil
	}

	return fmt.Errorf("health check against %s failed: %w", addr, err)
}

// KnownDependencyNames are the named grpc.health.v1 entries that services
// may publish for their dependencies (via UkamaGrpcServer.RegisterDependency).
// CheckServiceHealthDetailed queries each; services that don't register a
// given dependency return NOT_FOUND for it, which is treated as "not
// applicable" — so this list is safe to query against every service.
var KnownDependencyNames = []string{"db", "msgclient", "rabbitmq"}

// HealthReport is the detailed outcome of a service health check.
type HealthReport struct {
	// Err is non-nil when the service is unreachable or its default health
	// status is not SERVING.
	Err error

	// DegradedDeps lists dependencies reporting NOT_SERVING (e.g.
	// "db: NOT_SERVING"), regardless of whether Err is set.
	DegradedDeps []string
}

// CheckServiceHealthDetailed performs the same check as CheckServiceHealth
// and additionally queries the service's named dependency statuses so
// callers (e.g. GET /status) can distinguish "up but degraded" from "down".
func CheckServiceHealthDetailed(ctx context.Context, addr string) *HealthReport {
	report := &HealthReport{}

	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		report.Err = fmt.Errorf("failed to create connection to %s: %w", addr, err)

		return report
	}
	defer func() {
		_ = conn.Close()
	}()

	client := healthpb.NewHealthClient(conn)

	resp, cerr := client.Check(ctx, &healthpb.HealthCheckRequest{})
	switch {
	case cerr == nil && resp.GetStatus() == healthpb.HealthCheckResponse_SERVING:
		// Healthy so far; still collect dependency detail below.
	case cerr == nil:
		report.Err = fmt.Errorf("service at %s reported status %s",
			addr, resp.GetStatus().String())
	default:
		report.Err = fmt.Errorf("health check against %s failed: %w", addr, cerr)

		return report
	}

	for _, name := range KnownDependencyNames {
		dresp, derr := client.Check(ctx, &healthpb.HealthCheckRequest{Service: name})
		if derr != nil {
			// NOT_FOUND: the service does not publish this dependency.
			continue
		}

		if dresp.GetStatus() != healthpb.HealthCheckResponse_SERVING {
			report.DegradedDeps = append(report.DegradedDeps,
				fmt.Sprintf("%s: %s", name, dresp.GetStatus().String()))
		}
	}

	return report
}
