/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package server

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ukama/ukama/systems/analytics/collector/pkg"
	"github.com/ukama/ukama/systems/analytics/collector/pkg/db"
	"github.com/ukama/ukama/systems/analytics/collector/pkg/refresh"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/analytics/collector/pb/gen"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
)

const (
	sourceAll = "all"

	familyBusiness = "business"
	familyCustomer = "customer"
	familyNetwork  = "network"
	familyAll      = "all"
)

type CollectorServer struct {
	pb.UnimplementedCollectorServiceServer
	orgName     string
	stateRepo   db.StateRepo
	rollupRepo  db.RollupRepo
	eventRepo   db.EventRepo
	refresher   *refresh.Refresher
	msgbus      mb.MsgBusServiceClient
	pushGateway string
}

func NewCollectorServer(orgName string, stateRepo db.StateRepo, rollupRepo db.RollupRepo,
	eventRepo db.EventRepo, refresher *refresh.Refresher, msgBus mb.MsgBusServiceClient,
	pushGateway string) *CollectorServer {
	return &CollectorServer{
		orgName:     orgName,
		stateRepo:   stateRepo,
		rollupRepo:  rollupRepo,
		eventRepo:   eventRepo,
		refresher:   refresher,
		msgbus:      msgBus,
		pushGateway: pushGateway,
	}
}

func (c *CollectorServer) Refresh(ctx context.Context, req *pb.RefreshRequest) (*pb.RefreshResponse, error) {
	log.Infof("Refreshing source %s", req.GetSource())

	var sources []string

	if req.GetSource() == sourceAll {
		sources = refresh.Sources
	} else {
		if !isValidSource(req.GetSource()) {
			return nil, status.Errorf(codes.InvalidArgument,
				"invalid source %q: must be one of registry|subscriber|dataplan|metrics|node|inventory|billing|all",
				req.GetSource())
		}

		sources = []string{req.GetSource()}
	}

	states := make([]*pb.SourceState, 0, len(sources))

	for _, source := range sources {
		state, err := c.refresher.Refresh(source)
		if err != nil && state == nil {
			return nil, status.Errorf(codes.Internal,
				"failed to refresh source %s: %v", source, err)
		}

		states = append(states, dbRefreshStateToPb(state))
	}

	return &pb.RefreshResponse{States: states}, nil
}

func (c *CollectorServer) GetRefreshState(ctx context.Context, req *pb.GetRefreshStateRequest) (*pb.GetRefreshStateResponse, error) {
	states, err := c.stateRepo.GetRefreshStates()
	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"failed to get refresh states: %v", err)
	}

	rollups, err := c.stateRepo.GetRollupStates()
	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"failed to get rollup states: %v", err)
	}

	resp := &pb.GetRefreshStateResponse{
		States:  make([]*pb.SourceState, 0, len(states)),
		Rollups: make([]*pb.RollupState, 0, len(rollups)),
	}

	for i := range states {
		resp.States = append(resp.States, dbRefreshStateToPb(&states[i]))
	}

	for i := range rollups {
		resp.Rollups = append(resp.Rollups, dbRollupStateToPb(&rollups[i]))
	}

	return resp, nil
}

func (c *CollectorServer) RebuildRollups(ctx context.Context, req *pb.RebuildRollupsRequest) (*pb.RebuildRollupsResponse, error) {
	log.Infof("Rebuilding rollups for family %s", req.GetFamily())

	to := time.Now().UTC()
	from := to.AddDate(0, 0, -30)

	if req.GetFrom() != nil {
		from = req.GetFrom().AsTime()
	}

	if req.GetTo() != nil {
		to = req.GetTo().AsTime()
	}

	type rebuild struct {
		name string
		fn   func(from, to time.Time) error
	}

	var rebuilds []rebuild

	family := req.GetFamily()

	if family == familyBusiness || family == familyAll {
		rebuilds = append(rebuilds,
			rebuild{"business_sales_daily", c.rollupRepo.RebuildSalesDaily},
			rebuild{"business_package_daily", c.rollupRepo.RebuildPackageDaily},
			rebuild{"business_billing_daily", c.rollupRepo.RebuildBillingDaily},
		)
	}

	if family == familyCustomer || family == familyAll {
		rebuilds = append(rebuilds,
			rebuild{"customer_usage_daily", c.rollupRepo.RebuildCustomerUsageDaily},
			rebuild{"customer_state_daily", c.rollupRepo.RebuildCustomerStateDaily},
		)
	}

	if family == familyNetwork || family == familyAll {
		rebuilds = append(rebuilds,
			rebuild{"alarm_daily", c.rollupRepo.RebuildAlarmDaily},
			rebuild{"metric_hourly", c.rollupRepo.RebuildMetricHourly},
		)
	}

	if len(rebuilds) == 0 {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid family %q: must be one of business|customer|network|all", family)
	}

	for _, rb := range rebuilds {
		if err := rb.fn(from, to); err != nil {
			return nil, status.Errorf(codes.Internal,
				"failed to rebuild rollup %s: %v", rb.name, err)
		}

		if err := c.stateRepo.SetRollupWatermark(rb.name, to); err != nil {
			return nil, status.Errorf(codes.Internal,
				"failed to set watermark for rollup %s: %v", rb.name, err)
		}
	}

	rollups, err := c.stateRepo.GetRollupStates()
	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"failed to get rollup states: %v", err)
	}

	resp := &pb.RebuildRollupsResponse{
		Rollups: make([]*pb.RollupState, 0, len(rollups)),
	}

	for i := range rollups {
		resp.Rollups = append(resp.Rollups, dbRollupStateToPb(&rollups[i]))
	}

	return resp, nil
}

func (c *CollectorServer) SeedDemo(ctx context.Context, req *pb.SeedDemoRequest) (*pb.SeedDemoResponse, error) {
	if !pkg.IsDebugMode {
		return nil, status.Error(codes.PermissionDenied,
			"SeedDemo is only available in debug mode")
	}

	log.Infof("Seeding demo data: sites=%d nodes=%d customers=%d days=%d",
		req.GetSites(), req.GetNodes(), req.GetCustomers(), req.GetDays())

	/* Demo seeding is intentionally minimal: it marks all rollups dirty so a
	subsequent RebuildRollups pass regenerates them from any seeded facts. */
	for _, rollup := range []string{"business_sales_daily", "business_package_daily",
		"business_billing_daily", "customer_usage_daily", "customer_state_daily",
		"alarm_daily", "metric_hourly"} {
		if err := c.stateRepo.MarkRollupDirty(rollup); err != nil {
			return nil, status.Errorf(codes.Internal,
				"failed to mark rollup %s dirty: %v", rollup, err)
		}
	}

	return &pb.SeedDemoResponse{
		Detail: "demo seed requested; rollups marked dirty",
	}, nil
}

func isValidSource(source string) bool {
	for _, s := range refresh.Sources {
		if s == source {
			return true
		}
	}

	return false
}

func dbRefreshStateToPb(state *db.RefreshState) *pb.SourceState {
	s := &pb.SourceState{
		Source: state.Source,
		Status: state.Status,
		Detail: state.Detail,
	}

	if !state.LastRunAt.IsZero() {
		s.LastRunAt = timestamppb.New(state.LastRunAt)
	}

	if !state.LastSuccessAt.IsZero() {
		s.LastSuccessAt = timestamppb.New(state.LastSuccessAt)
	}

	return s
}

func dbRollupStateToPb(state *db.RollupState) *pb.RollupState {
	s := &pb.RollupState{
		Rollup: state.Rollup,
		Dirty:  state.Dirty,
	}

	if !state.Watermark.IsZero() {
		s.Watermark = timestamppb.New(state.Watermark)
	}

	return s
}
