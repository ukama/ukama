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
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ukama/ukama/systems/analytics/customer/pkg"
	"github.com/ukama/ukama/systems/analytics/customer/pkg/db"
	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/uuid"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/analytics/customer/pb/gen"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
)

const uuidParsingError = "Error parsing UUID"

type CustomerServer struct {
	pb.UnimplementedCustomerServiceServer
	orgName              string
	customerRepo         db.CustomerRepo
	simRepo              db.SimRepo
	supportRepo          db.SupportRepo
	msgbus               mb.MsgBusServiceClient
	baseRoutingKey       msgbus.RoutingKeyBuilder
	simLowStockThreshold uint32
}

func NewCustomerServer(orgName string, customerRepo db.CustomerRepo, simRepo db.SimRepo,
	supportRepo db.SupportRepo, msgBus mb.MsgBusServiceClient, simLowStockThreshold uint32) *CustomerServer {
	return &CustomerServer{
		orgName:              orgName,
		customerRepo:         customerRepo,
		simRepo:              simRepo,
		supportRepo:          supportRepo,
		msgbus:               msgBus,
		baseRoutingKey:       msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		simLowStockThreshold: simLowStockThreshold,
	}
}

/* GetOverview returns customer base KPIs for the requested window with
deltas against the previous window of equal length. */
func (c *CustomerServer) GetOverview(ctx context.Context, req *pb.GetOverviewRequest) (*pb.GetOverviewResponse, error) {
	log.Infof("GetOverview: %v", req)

	networkId, err := parseOptionalUuid(req.GetNetworkId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"%s of network. Error %s", uuidParsingError, err.Error())
	}

	w, err := resolveWindow(req.GetWindow(), time.Now())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid window: %s", err.Error())
	}

	total, active, newCount, expired, failed, err := c.customerRepo.Counts(networkId, w.From, w.To)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "customers")
	}

	_, _, prevNew, _, prevFailed, err := c.customerRepo.Counts(networkId, w.PrevFrom, w.PrevTo)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "customers")
	}

	asOf := timestamppb.Now()

	kpis := []*pb.Kpi{
		countKpi("total_customers", total, 0, "", asOf),
		countKpi("active_customers", active, 0, "", asOf),
		countKpi("new_customers", newCount, float64(newCount)-float64(prevNew), w.Period, asOf),
		countKpi("expired_customers", expired, 0, "", asOf),
		countKpi("failed_activations", failed, float64(failed)-float64(prevFailed), w.Period, asOf),
	}

	return &pb.GetOverviewResponse{
		Kpis: kpis,
	}, nil
}

/* List returns a page of customers filtered by network, site and status. */
func (c *CustomerServer) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	log.Infof("List: %v", req)

	networkId, err := parseOptionalUuid(req.GetNetworkId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"%s of network. Error %s", uuidParsingError, err.Error())
	}

	siteId, err := parseOptionalUuid(req.GetSiteId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"%s of site. Error %s", uuidParsingError, err.Error())
	}

	page, pageSize := normalizePagination(req.GetPage(), req.GetPageSize())

	customers, count, err := c.customerRepo.List(networkId, siteId, req.GetStatus(), page, pageSize)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "customers")
	}

	rows, err := c.toCustomerRows(customers)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sites")
	}

	return &pb.ListResponse{
		Customers: rows,
		Meta:      buildMeta(count, page, pageSize),
	}, nil
}

/* Search performs a case-insensitive match against name, email and sim iccid. */
func (c *CustomerServer) Search(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
	log.Infof("Search: %v", req)

	if req.GetQuery() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "query cannot be empty")
	}

	networkId, err := parseOptionalUuid(req.GetNetworkId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"%s of network. Error %s", uuidParsingError, err.Error())
	}

	page, pageSize := normalizePagination(req.GetPage(), req.GetPageSize())

	customers, count, err := c.customerRepo.Search(req.GetQuery(), networkId, page, pageSize)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "customers")
	}

	rows, err := c.toCustomerRows(customers)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sites")
	}

	return &pb.SearchResponse{
		Customers: rows,
		Meta:      buildMeta(count, page, pageSize),
	}, nil
}

/* Get returns one customer with detail KPIs and package history. */
func (c *CustomerServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	log.Infof("Get: %v", req)

	customerId, err := uuid.FromString(req.GetCustomerId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"%s of customer. Error %s", uuidParsingError, err.Error())
	}

	w, err := resolveWindow(req.GetWindow(), time.Now())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid window: %s", err.Error())
	}

	customer, err := c.customerRepo.Get(customerId)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "customer")
	}

	rows, err := c.toCustomerRows([]db.CustomerSnapshot{*customer})
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sites")
	}

	usage, err := c.customerRepo.UsageBetween(customerId, w.From, w.To)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "usage")
	}

	prevUsage, err := c.customerRepo.UsageBetween(customerId, w.PrevFrom, w.PrevTo)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "usage")
	}

	intervals, err := c.customerRepo.PackageIntervals(customerId)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "package intervals")
	}

	rows[0].DataUsage = usage

	asOf := timestamppb.Now()

	lastSeenMinutes := float64(0)
	if customer.LastSeenAt != nil {
		lastSeenMinutes = time.Since(*customer.LastSeenAt).Minutes()
	}

	kpis := []*pb.Kpi{
		{
			Key:         "data_usage",
			Value:       usage,
			Formatted:   fmt.Sprintf("%.1f MB", usage),
			Delta:       usage - prevUsage,
			DeltaPeriod: w.Period,
			AsOf:        asOf,
		},
		{
			Key:       "last_seen_minutes",
			Value:     lastSeenMinutes,
			Formatted: fmt.Sprintf("%.0f min", lastSeenMinutes),
			AsOf:      asOf,
		},
	}

	history := make([]*pb.PackageInterval, 0, len(intervals))
	for _, iv := range intervals {
		p := &pb.PackageInterval{
			PackageId: iv.PackageId.String(),
			State:     iv.State,
			StartAt:   timestamppb.New(iv.StartAt),
		}

		if iv.PackageId == customer.PackageId {
			p.PackageName = customer.PackageName
		}

		if iv.EndAt != nil {
			p.EndAt = timestamppb.New(*iv.EndAt)
		}

		history = append(history, p)
	}

	return &pb.GetResponse{
		Customer:       rows[0],
		Kpis:           kpis,
		PackageHistory: history,
	}, nil
}

/* GetSupport returns the support diagnosis for a customer: derived signals,
likely issue, recommended action and recent activity. */
func (c *CustomerServer) GetSupport(ctx context.Context, req *pb.GetSupportRequest) (*pb.GetSupportResponse, error) {
	log.Infof("GetSupport: %v", req)

	customerId, err := uuid.FromString(req.GetCustomerId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"%s of customer. Error %s", uuidParsingError, err.Error())
	}

	customer, err := c.customerRepo.Get(customerId)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "customer")
	}

	rows, err := c.toCustomerRows([]db.CustomerSnapshot{*customer})
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sites")
	}

	now := time.Now()

	in := diagnosisInput{
		Now:           now,
		SimStatus:     customer.SimStatus,
		PackageStatus: customer.PackageStatus,
		LastSeenAt:    customer.LastSeenAt,
	}

	if customer.SiteId != uuid.Nil {
		health, err := c.supportRepo.SiteHealth(customer.SiteId)
		if err == nil && health != nil {
			in.HasSiteHealth = true
			in.SiteUptimePercent = health.UptimePercent
		} else if err != nil {
			// missing site health is a soft failure: signal stays unknown
			log.Warnf("no site health for site %s: %v", customer.SiteId, err)
		}
	}

	usage24h, err := c.customerRepo.UsageBetween(customerId, now.Add(-24*time.Hour), now)
	if err != nil {
		log.Warnf("failed reading usage for customer %s: %v", customerId, err)
	} else {
		in.UsageLast24hMb = usage24h
	}

	res := diagnose(in)

	signals := make([]*pb.SupportSignal, 0, len(res.Signals))
	for _, s := range res.Signals {
		signals = append(signals, &pb.SupportSignal{
			Key:    s.Key,
			State:  s.State,
			Detail: s.Detail,
		})
	}

	logs, err := c.supportRepo.RecentActivityFor(customerId, 10)
	if err != nil {
		log.Warnf("failed reading recent activity for customer %s: %v", customerId, err)
	}

	activity := make([]*pb.ActivityItem, 0, len(logs))
	for _, l := range logs {
		activity = append(activity, &pb.ActivityItem{
			RoutingKey:  l.RoutingKey,
			Description: string(l.Payload),
			OccurredAt:  timestamppb.New(l.OccurredAt),
		})
	}

	return &pb.GetSupportResponse{
		Customer:          rows[0],
		LikelyIssue:       res.LikelyIssue,
		RecommendedAction: res.RecommendedAction,
		EscalationNeeded:  res.EscalationNeeded,
		Signals:           signals,
		RecentActivity:    activity,
	}, nil
}

/* GetSims returns a page of sims filtered by network and status. */
func (c *CustomerServer) GetSims(ctx context.Context, req *pb.GetSimsRequest) (*pb.GetSimsResponse, error) {
	log.Infof("GetSims: %v", req)

	networkId, err := parseOptionalUuid(req.GetNetworkId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"%s of network. Error %s", uuidParsingError, err.Error())
	}

	page, pageSize := normalizePagination(req.GetPage(), req.GetPageSize())

	sims, count, err := c.simRepo.List(networkId, req.GetStatus(), page, pageSize)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sims")
	}

	rows := make([]*pb.SimRow, 0, len(sims))
	for _, s := range sims {
		row := &pb.SimRow{
			SimId:   s.SimId,
			Iccid:   s.Iccid,
			Status:  s.Status,
			BatchId: s.BatchId,
		}

		if s.CustomerId != uuid.Nil {
			row.CustomerId = s.CustomerId.String()
		}

		if s.AllocatedAt != nil {
			row.AllocatedAt = timestamppb.New(*s.AllocatedAt)
		}

		rows = append(rows, row)
	}

	return &pb.GetSimsResponse{
		Sims: rows,
		Meta: buildMeta(count, page, pageSize),
	}, nil
}

/* GetSimPool returns sim pool KPIs (including low stock) and batches. */
func (c *CustomerServer) GetSimPool(ctx context.Context, req *pb.GetSimPoolRequest) (*pb.GetSimPoolResponse, error) {
	log.Infof("GetSimPool: %v", req)

	total, available, active, assigned, suspended, faulty, err := c.simRepo.PoolCounts()
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sims")
	}

	batches, err := c.simRepo.Batches()
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sim batches")
	}

	asOf := timestamppb.Now()

	lowStock := float64(0)
	lowStockFormatted := "ok"
	if available < c.simLowStockThreshold {
		lowStock = 1
		lowStockFormatted = fmt.Sprintf("low: %d available (threshold %d)",
			available, c.simLowStockThreshold)
	}

	kpis := []*pb.Kpi{
		countKpi("total_sims", total, 0, "", asOf),
		countKpi("available_sims", available, 0, "", asOf),
		countKpi("active_sims", active, 0, "", asOf),
		countKpi("assigned_sims", assigned, 0, "", asOf),
		countKpi("suspended_sims", suspended, 0, "", asOf),
		countKpi("faulty_sims", faulty, 0, "", asOf),
		{
			Key:       "low_stock",
			Value:     lowStock,
			Formatted: lowStockFormatted,
			AsOf:      asOf,
		},
	}

	pbBatches := make([]*pb.SimBatch, 0, len(batches))
	for _, b := range batches {
		assignedPercent := float64(0)
		if b.Quantity > 0 {
			assignedPercent = float64(b.Assigned) / float64(b.Quantity) * 100
		}

		batch := &pb.SimBatch{
			BatchId:         b.BatchId,
			Quantity:        b.Quantity,
			Assigned:        b.Assigned,
			AssignedPercent: assignedPercent,
		}

		if b.UploadedAt != nil {
			batch.UploadedAt = timestamppb.New(*b.UploadedAt)
		}

		pbBatches = append(pbBatches, batch)
	}

	return &pb.GetSimPoolResponse{
		Kpis:    kpis,
		Batches: pbBatches,
	}, nil
}

/* helpers */

// parseOptionalUuid parses an optional uuid filter; empty means "no filter".
func parseOptionalUuid(s string) (uuid.UUID, error) {
	if s == "" {
		return uuid.Nil, nil
	}

	return uuid.FromString(s)
}

func normalizePagination(page, pageSize uint32) (uint32, uint32) {
	if page == 0 {
		page = 1
	}

	if pageSize == 0 {
		pageSize = db.DefaultPageSize
	}

	return page, pageSize
}

func buildMeta(count int64, page, pageSize uint32) *pb.Meta {
	pages := uint32(0)
	if pageSize > 0 {
		pages = uint32((count + int64(pageSize) - 1) / int64(pageSize))
	}

	return &pb.Meta{
		Count: uint32(count),
		Page:  page,
		Size:  pageSize,
		Pages: pages,
	}
}

func countKpi(key string, value uint32, delta float64, deltaPeriod string, asOf *timestamppb.Timestamp) *pb.Kpi {
	return &pb.Kpi{
		Key:         key,
		Value:       float64(value),
		Formatted:   fmt.Sprintf("%d", value),
		Delta:       delta,
		DeltaPeriod: deltaPeriod,
		AsOf:        asOf,
	}
}

// toCustomerRows maps snapshots to pb rows, resolving site names in one query.
func (c *CustomerServer) toCustomerRows(customers []db.CustomerSnapshot) ([]*pb.CustomerRow, error) {
	siteIds := make([]uuid.UUID, 0, len(customers))
	seen := make(map[uuid.UUID]bool)

	for _, cu := range customers {
		if cu.SiteId != uuid.Nil && !seen[cu.SiteId] {
			seen[cu.SiteId] = true
			siteIds = append(siteIds, cu.SiteId)
		}
	}

	siteNames, err := c.customerRepo.SiteNames(siteIds)
	if err != nil {
		return nil, err
	}

	rows := make([]*pb.CustomerRow, 0, len(customers))
	for _, cu := range customers {
		row := &pb.CustomerRow{
			CustomerId:    cu.CustomerId.String(),
			Name:          cu.Name,
			Email:         cu.Email,
			Status:        cu.Status,
			PackageName:   cu.PackageName,
			PackageStatus: cu.PackageStatus,
			SimIccid:      cu.SimIccid,
			SimStatus:     cu.SimStatus,
		}

		if cu.SiteId != uuid.Nil {
			row.SiteId = cu.SiteId.String()
			row.SiteName = siteNames[cu.SiteId]
		}

		if cu.LastSeenAt != nil {
			row.LastSeen = timestamppb.New(*cu.LastSeenAt)
		}

		rows = append(rows, row)
	}

	return rows, nil
}
