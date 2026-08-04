/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

import (
	"context"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/wI2L/fizz"

	ugrpc "github.com/ukama/ukama/systems/common/grpc"
)

const (
	ServiceStatusAvailable   = "available"
	ServiceStatusUnavailable = "unavailable"

	SystemStatusOk       = "ok"
	SystemStatusDegraded = "degraded"
	SystemStatusDown     = "down"

	defaultStatusCheckTimeout = 5 * time.Second
)

type ServiceStatus struct {
	Name   string `json:"name"`
	Host   string `json:"host"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type StatusResponse struct {
	System   string          `json:"system"`
	Status   string          `json:"status"`
	Services []ServiceStatus `json:"services"`
}

// RegisterStatusEndpoint adds a GET /status route to the given fizz router.
// It health-checks every gRPC service of the system (name -> host:port) in
// parallel via grpc.health.v1 (with fallback to the legacy ukama.health.v1)
// and reports each one as available/unavailable, plus an aggregated system
// status: ok (all up), degraded (some up), or down (all down).
func RegisterStatusEndpoint(f *fizz.Fizz, systemName string,
	services map[string]string, timeout time.Duration) {
	if timeout <= 0 {
		timeout = defaultStatusCheckTimeout
	}

	// fizz derives OpenAPI operation IDs from handler function names;
	// anonymous handlers (like the /ping one) all reduce to "func1", so an
	// explicit ID is required to avoid a duplicate-ID panic at startup.
	f.GET("/status", []fizz.OperationOption{
		fizz.ID("getSystemStatus"),
		fizz.Summary("Get system status"),
		fizz.Description("Health-checks all gRPC services of this system and reports each as available/unavailable."),
	}, tonic.Handler(func(c *gin.Context) (*StatusResponse, error) {
		return checkServices(c.Request.Context(), systemName, services, timeout), nil
	}, http.StatusOK))
}

func checkServices(ctx context.Context, systemName string,
	services map[string]string, timeout time.Duration) *StatusResponse {
	var (
		mu      sync.Mutex
		wg      sync.WaitGroup
		results = make([]ServiceStatus, 0, len(services))
	)

	for name, host := range services {
		wg.Add(1)

		go func(name, host string) {
			defer wg.Done()

			cctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			s := ServiceStatus{
				Name:   name,
				Host:   host,
				Status: ServiceStatusAvailable,
			}

			if err := ugrpc.CheckServiceHealth(cctx, host); err != nil {
				s.Status = ServiceStatusUnavailable
				s.Error = err.Error()
			}

			mu.Lock()
			results = append(results, s)
			mu.Unlock()
		}(name, host)
	}

	wg.Wait()

	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})

	unavailable := 0
	for _, s := range results {
		if s.Status == ServiceStatusUnavailable {
			unavailable++
		}
	}

	systemStatus := SystemStatusOk
	switch {
	case len(results) > 0 && unavailable == len(results):
		systemStatus = SystemStatusDown
	case unavailable > 0:
		systemStatus = SystemStatusDegraded
	}

	return &StatusResponse{
		System:   systemName,
		Status:   systemStatus,
		Services: results,
	}
}
