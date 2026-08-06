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

// LivenessService is the grpc.health.v1 service name for Kubernetes
// liveness probes (grpc: {port: 9090, service: "live"}). It reports SERVING
// while the process is healthy and, unlike the default "" service, does NOT
// react to short dependency outages — only readiness does. It escalates to
// NOT_SERVING when a critical dependency stays down for
// DependencyRestartThreshold consecutive checks, so kubelet restarts a
// container that has been degraded for a prolonged period (fresh
// connections often fix wedged pools/consumers), while brief blips never
// cause restarts.
const LivenessService = "live"

// All dependency state transitions are driven by ATTEMPT COUNTS, never by
// elapsed time: each tick every check is attempted once, a failed attempt
// increments that dependency's consecutive-failure counter, and a
// successful one resets it to zero. The thresholds below are counts of
// consecutive failed attempts. Durations in comments are only the derived
// worst case (count x interval) for Kubernetes probe tuning.
const (
	// defaultDependencyCheckInterval is how often one check attempt runs
	// for each registered dependency. This is the cadence the failure
	// counters advance at; it is NOT a health criterion by itself.
	defaultDependencyCheckInterval = 15 * time.Second

	// defaultDependencyCheckTimeout is the deadline embedded in the
	// context of ONE check attempt. It bounds a single call so a hung
	// dependency cannot stall the monitor; it is NOT a health criterion.
	// An attempt that exceeds it simply counts as one more failure.
	defaultDependencyCheckTimeout = 3 * time.Second

	// defaultDependencyFailThreshold is the count of consecutive failed
	// attempts after which a dependency is marked down: its named status
	// and, when critical, the default "" (readiness) status flip to
	// NOT_SERVING (3 failures ~= 45s at the 15s interval).
	defaultDependencyFailThreshold = 3

	// defaultDependencyRestartThreshold is the count of consecutive
	// failed attempts of a critical dependency after which the "live"
	// status also flips to NOT_SERVING, asking kubelet to restart the
	// container (4 failures ~= 60s at the 15s interval). Must be greater
	// than the fail threshold so readiness always fails before liveness.
	defaultDependencyRestartThreshold = 4
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

	// DependencyCheckInterval is how often one check attempt runs for
	// each registered dependency.
	DependencyCheckInterval time.Duration

	// DependencyCheckTimeout is the per-attempt call deadline (embedded
	// in each check's context). It is not a health criterion.
	DependencyCheckTimeout time.Duration

	// DependencyFailThreshold is the count of consecutive failed
	// attempts after which a dependency is marked down (readiness).
	DependencyFailThreshold int

	// DependencyRestartThreshold is the count of consecutive failed
	// attempts of a critical dependency after which the LivenessService
	// ("live") status flips to NOT_SERVING so kubelet restarts the
	// container. Zero uses defaultDependencyRestartThreshold. Values <=
	// the fail threshold are raised to failThreshold+1 so readiness
	// always fails before liveness does.
	DependencyRestartThreshold int
}

func (g *UkamaGrpcServer) startDependencyMonitor() {
	g.depMu.Lock()
	n := len(g.deps)
	g.depMu.Unlock()

	if n == 0 {
		return
	}

	interval := g.DependencyCheckInterval
	if interval <= 0 {
		interval = defaultDependencyCheckInterval
	}
	checkTimeout := g.DependencyCheckTimeout
	if checkTimeout <= 0 {
		checkTimeout = defaultDependencyCheckTimeout
	}
	threshold := g.DependencyFailThreshold
	if threshold <= 0 {
		threshold = defaultDependencyFailThreshold
	}
	restartThreshold := g.DependencyRestartThreshold
	if restartThreshold <= 0 {
		restartThreshold = defaultDependencyRestartThreshold
	}
	if restartThreshold <= threshold {
		restartThreshold = threshold + 1
	}

	g.depStop = make(chan struct{})

	log.Infof("Starting dependency health monitor: %d checks, every %s", n, interval)

	go func() {
		// Run once immediately so status is meaningful right after boot.
		g.runDependencyChecks(checkTimeout, threshold, restartThreshold)

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				g.runDependencyChecks(checkTimeout, threshold, restartThreshold)
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

func (g *UkamaGrpcServer) runDependencyChecks(checkTimeout time.Duration, threshold, restartThreshold int) {
	g.depMu.Lock()
	defer g.depMu.Unlock()

	criticalDown := false
	criticalRestart := false

	for _, d := range g.deps {
		ctx, cancel := context.WithTimeout(context.Background(), checkTimeout)
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
			if d.critical && d.failures == restartThreshold {
				log.Errorf("Critical dependency %q still down after %d consecutive failures: "+
					"flipping liveness (%q) to NOT_SERVING to request a container restart",
					d.name, d.failures, LivenessService)
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
				if d.failures >= restartThreshold {
					criticalRestart = true
				}
			}
		}
		g.GrpcHealth.SetServingStatus(d.name, st)
	}

	overall := healthpb.HealthCheckResponse_SERVING
	if criticalDown {
		overall = healthpb.HealthCheckResponse_NOT_SERVING
	}
	g.GrpcHealth.SetServingStatus("", overall)

	// Escalation: a critical dependency down for restartThreshold
	// consecutive checks flips liveness so kubelet recreates the
	// container; recovery flips it back.
	live := healthpb.HealthCheckResponse_SERVING
	if criticalRestart {
		live = healthpb.HealthCheckResponse_NOT_SERVING
	}
	g.GrpcHealth.SetServingStatus(LivenessService, live)
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
