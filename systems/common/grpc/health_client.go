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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"

	upb "github.com/ukama/ukama/systems/common/pb/gen/health"
)

// CheckServiceHealth dials addr and performs a health check against the
// standard grpc.health.v1.Health service. If the target does not implement
// it (e.g. a service not yet rebuilt with the standard health server), it
// falls back to the legacy ukama.health.v1.Health service.
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

	if status.Code(err) != codes.Unimplemented {
		return fmt.Errorf("health check against %s failed: %w", addr, err)
	}

	// Fallback: legacy ukama.health.v1 health service.
	lresp, lerr := upb.NewHealthClient(conn).Check(ctx, &upb.HealthCheckRequest{})
	if lerr != nil {
		return fmt.Errorf("legacy health check against %s failed: %w", addr, lerr)
	}

	if lresp.GetStatus() != upb.HealthCheckResponse_SERVING {
		return fmt.Errorf("service at %s reported legacy status %s",
			addr, lresp.GetStatus().String())
	}

	return nil
}
