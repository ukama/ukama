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
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	usql "github.com/ukama/ukama/systems/common/sql"
)

// LivenessService is the grpc.health.v1 service name that always reports
// SERVING while the process is up. Point Kubernetes liveness probes at it
// (grpc: {port: 9090, service: "live"}) so dependency failures never cause
// pod restarts; only readiness (the default "" service) reflects them.
const LivenessService = "live"

const (
	defaultDependencyInterval      = 15 * time.Second
	defaultDependencyTimeout       = 3 * time.Second
	defaultDependencyFailThreshold = 3
)

// dependency is one registered dependency check with its flap state.
type dependency struct {
	name     string
	critical bool
	check    func(ctx context.Context) error

	failures int
	down     bool
	lastErr  error
}

// RegisterDependency adds a named dependency check (e.g. "db", "msgclient",
// "rabbitmq") that the server monitors in the background once started.
//
// The result of each check is published on the standard health server under
// the dependency's name. When any *critical* dependency is down, the default
// ("") health service flips to NOT_SERVING so readiness probes and /status
// see it; non-critical dependencies only affect their own named status.
//
// Registration is opt-in: services that register nothing behave exactly as
// before (no monitor runs, default status stays SERVING).
// Must be called before StartServer.
func (g *UkamaGrpcServer) RegisterDependency(name string, critical bool,
	check func(ctx context.Context) error) {
	g.depMu.Lock()
	defer g.depMu.Unlock()

	g.deps = append(g.deps, &dependency{name: name, critical: critical, check: check})
}

// dependencyState holds the monitor plumbing embedded in UkamaGrpcServer.
type dependencyState struct {
	depMu   sync.Mutex
	deps    []*dependency
	depStop chan struct{}

	// Overridable before StartServer; zero values use defaults.
	DependencyInterval      time.Duration
	DependencyTimeout       time.Duration
	DependencyFailThreshold int
}

func (g *UkamaGrpcServer) startDependencyMonitor() {
	g.depMu.Lock()
	n := len(g.deps)
	g.depMu.Unlock()

	if n == 0 {
		return
	}

	interval := g.DependencyInterval
	if interval <= 0 {
		interval = defaultDependencyInterval
	}
	timeout := g.DependencyTimeout
	if timeout <= 0 {
		timeout = defaultDependencyTimeout
	}
	threshold := g.DependencyFailThreshold
	if threshold <= 0 {
		threshold = defaultDependencyFailThreshold
	}

	g.depStop = make(chan struct{})

	log.Infof("Starting dependency health monitor: %d checks, every %s", n, interval)

	go func() {
		// Run once immediately so status is meaningful right after boot.
		g.runDependencyChecks(timeout, threshold)

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				g.runDependencyChecks(timeout, threshold)
			case <-g.depStop:
				return
			}
		}
	}()
}

func (g *UkamaGrpcServer) stopDependencyMonitor() {
	if g.depStop != nil {
		close(g.depStop)
		g.depStop = nil
	}
}

func (g *UkamaGrpcServer) runDependencyChecks(timeout time.Duration, threshold int) {
	g.depMu.Lock()
	defer g.depMu.Unlock()

	criticalDown := false

	for _, d := range g.deps {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		err := d.check(ctx)
		cancel()

		if err != nil {
			d.failures++
			d.lastErr = err
			if d.failures >= threshold && !d.down {
				d.down = true
				log.Errorf("Dependency %q is DOWN after %d consecutive failures: %v",
					d.name, d.failures, err)
			}
		} else {
			if d.down {
				log.Infof("Dependency %q recovered", d.name)
			}
			d.failures = 0
			d.down = false
			d.lastErr = nil
		}

		st := healthpb.HealthCheckResponse_SERVING
		if d.down {
			st = healthpb.HealthCheckResponse_NOT_SERVING
			if d.critical {
				criticalDown = true
			}
		}
		g.GrpcHealth.SetServingStatus(d.name, st)
	}

	overall := healthpb.HealthCheckResponse_SERVING
	if criticalDown {
		overall = healthpb.HealthCheckResponse_NOT_SERVING
	}
	g.GrpcHealth.SetServingStatus("", overall)
}

// DBCheck returns a dependency check that pings the service's database.
func DBCheck(db usql.Db) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		sqlDB, err := db.GetGormDb().DB()
		if err != nil {
			return err
		}

		return sqlDB.PingContext(ctx)
	}
}

// MsgClientCheck returns a dependency check that health-checks the system's
// msgclient service (grpc.health.v1 with legacy fallback). Because msgclient
// registers its own "rabbitmq" dependency, a broken message bus propagates
// transitively to every service using this check.
func MsgClientCheck(msgclientHost string) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		return CheckServiceHealth(ctx, msgclientHost)
	}
}

// AmqpCheck returns a dependency check that verifies the RabbitMQ broker at
// uri accepts AMQP connections. It dials a fresh short-lived connection per
// check, which keeps it independent of the consumer/publisher state.
func AmqpCheck(uri string) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		type result struct{ err error }
		done := make(chan result, 1)

		go func() {
			conn, err := amqp.Dial(uri)
			if err == nil {
				_ = conn.Close()
			}
			done <- result{err: err}
		}()

		select {
		case r := <-done:
			return r.err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
