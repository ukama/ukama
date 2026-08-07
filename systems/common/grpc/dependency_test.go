/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package grpc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// Mirrors the production defaults: readiness at 3 consecutive failures
// (~45s at the 15s interval), liveness at 4 (~60s).
const (
	testFailThreshold    = defaultDependencyFailThreshold
	testRestartThreshold = defaultDependencyRestartThreshold
)

// newTestServer returns an UkamaGrpcServer with a health server, one
// registered dependency driven by the returned error pointer, and both
// statuses initialized to SERVING (as startServerInternal does).
func newTestServer(name string, critical bool, depErr *error) *UkamaGrpcServer {
	g := &UkamaGrpcServer{GrpcHealth: health.NewServer()}

	g.GrpcHealth.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	g.GrpcHealth.SetServingStatus(LivenessService, healthpb.HealthCheckResponse_SERVING)

	g.RegisterDependency(name, critical, func(ctx context.Context) error {
		return *depErr
	})

	return g
}

func checkStatus(t *testing.T, g *UkamaGrpcServer, service string) healthpb.HealthCheckResponse_ServingStatus {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := g.GrpcHealth.Check(ctx, &healthpb.HealthCheckRequest{Service: service})
	assert.NoError(t, err)

	return resp.GetStatus()
}

func runChecks(g *UkamaGrpcServer, n int) {
	for i := 0; i < n; i++ {
		g.runDependencyChecks(time.Second, testFailThreshold, testRestartThreshold)
	}
}

func TestDependencyEscalation_CriticalDep(t *testing.T) {
	depErr := error(nil)
	g := newTestServer("db", true, &depErr)

	// Healthy: everything SERVING.
	runChecks(g, 1)
	assert.Equal(t, healthpb.HealthCheckResponse_SERVING, checkStatus(t, g, ""))
	assert.Equal(t, healthpb.HealthCheckResponse_SERVING, checkStatus(t, g, LivenessService))
	assert.Equal(t, healthpb.HealthCheckResponse_SERVING, checkStatus(t, g, "db"))

	depErr = errors.New("connection refused")

	// Below fail threshold (2 failures): blip absorbed, nothing flips.
	runChecks(g, testFailThreshold-1)
	assert.Equal(t, healthpb.HealthCheckResponse_SERVING, checkStatus(t, g, ""))
	assert.Equal(t, healthpb.HealthCheckResponse_SERVING, checkStatus(t, g, LivenessService))

	// At fail threshold (3rd failure): readiness ("") fails, liveness holds.
	runChecks(g, 1)
	assert.Equal(t, healthpb.HealthCheckResponse_NOT_SERVING, checkStatus(t, g, ""))
	assert.Equal(t, healthpb.HealthCheckResponse_NOT_SERVING, checkStatus(t, g, "db"))
	assert.Equal(t, healthpb.HealthCheckResponse_SERVING, checkStatus(t, g, LivenessService))

	// Just below restart threshold: liveness still holds.
	runChecks(g, testRestartThreshold-testFailThreshold-1)
	assert.Equal(t, healthpb.HealthCheckResponse_SERVING, checkStatus(t, g, LivenessService))

	// At restart threshold: liveness fails -> kubelet restarts.
	runChecks(g, 1)
	assert.Equal(t, healthpb.HealthCheckResponse_NOT_SERVING, checkStatus(t, g, LivenessService))
	assert.Equal(t, healthpb.HealthCheckResponse_NOT_SERVING, checkStatus(t, g, ""))

	// Recovery: one successful check restores everything.
	depErr = nil
	runChecks(g, 1)
	assert.Equal(t, healthpb.HealthCheckResponse_SERVING, checkStatus(t, g, ""))
	assert.Equal(t, healthpb.HealthCheckResponse_SERVING, checkStatus(t, g, LivenessService))
	assert.Equal(t, healthpb.HealthCheckResponse_SERVING, checkStatus(t, g, "db"))
}

func TestDependencyEscalation_NonCriticalDepNeverEscalates(t *testing.T) {
	depErr := errors.New("smtp down")
	g := newTestServer("mailer", false, &depErr)

	// Way past both thresholds: only the named status flips.
	runChecks(g, testRestartThreshold*2)
	assert.Equal(t, healthpb.HealthCheckResponse_NOT_SERVING, checkStatus(t, g, "mailer"))
	assert.Equal(t, healthpb.HealthCheckResponse_SERVING, checkStatus(t, g, ""))
	assert.Equal(t, healthpb.HealthCheckResponse_SERVING, checkStatus(t, g, LivenessService))
}

func TestDependencyEscalation_RecoveryBetweenThresholds(t *testing.T) {
	depErr := errors.New("timeout")
	g := newTestServer("msgclient", true, &depErr)

	// Fail past readiness but short of restart, then recover: the failure
	// counter resets, so a new outage starts the escalation from zero.
	runChecks(g, testRestartThreshold-1)
	assert.Equal(t, healthpb.HealthCheckResponse_NOT_SERVING, checkStatus(t, g, ""))
	assert.Equal(t, healthpb.HealthCheckResponse_SERVING, checkStatus(t, g, LivenessService))

	depErr = nil
	runChecks(g, 1)
	assert.Equal(t, healthpb.HealthCheckResponse_SERVING, checkStatus(t, g, ""))

	depErr = errors.New("timeout again")
	runChecks(g, testFailThreshold-1)
	assert.Equal(t, healthpb.HealthCheckResponse_SERVING, checkStatus(t, g, ""),
		"failure counter must reset after recovery")
	assert.Equal(t, healthpb.HealthCheckResponse_SERVING, checkStatus(t, g, LivenessService))
}

func TestDependencyRestartThreshold_ClampedAboveFailThreshold(t *testing.T) {
	depErr := errors.New("down")
	g := newTestServer("db", true, &depErr)

	// restartThreshold <= failThreshold must be raised to failThreshold+1
	// by startDependencyMonitor; emulate its clamping here and verify
	// readiness always fails strictly before liveness.
	g.DependencyFailThreshold = 3
	g.DependencyRestartThreshold = 2

	threshold := g.DependencyFailThreshold
	restart := g.DependencyRestartThreshold
	if restart <= threshold {
		restart = threshold + 1
	}

	for i := 0; i < threshold; i++ {
		g.runDependencyChecks(time.Second, threshold, restart)
	}
	assert.Equal(t, healthpb.HealthCheckResponse_NOT_SERVING, checkStatus(t, g, ""))
	assert.Equal(t, healthpb.HealthCheckResponse_SERVING, checkStatus(t, g, LivenessService))

	g.runDependencyChecks(time.Second, threshold, restart)
	assert.Equal(t, healthpb.HealthCheckResponse_NOT_SERVING, checkStatus(t, g, LivenessService))
}
