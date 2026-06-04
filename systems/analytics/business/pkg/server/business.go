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

	"github.com/ukama/ukama/systems/common/grpc"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/analytics/business/pb/gen"
	"github.com/ukama/ukama/systems/analytics/business/pkg/db"
)

const defaultPageSize = 20

// BusinessServer serves read-only business analytics KPIs computed from the
// shared analytics database (rollups, snapshots and fact tables maintained
// by the collector service).
type BusinessServer struct {
	pb.UnimplementedBusinessServiceServer
	orgName       string
	salesRepo     db.SalesRepo
	packageRepo   db.PackageRepo
	siteRepo      db.SiteRepo
	billingRepo   db.BillingRepo
	inventoryRepo db.InventoryRepo
	activityRepo  db.ActivityRepo
}

func NewBusinessServer(orgName string, salesRepo db.SalesRepo, packageRepo db.PackageRepo,
	siteRepo db.SiteRepo, billingRepo db.BillingRepo, inventoryRepo db.InventoryRepo,
	activityRepo db.ActivityRepo) *BusinessServer {
	return &BusinessServer{
		orgName:       orgName,
		salesRepo:     salesRepo,
		packageRepo:   packageRepo,
		siteRepo:      siteRepo,
		billingRepo:   billingRepo,
		inventoryRepo: inventoryRepo,
		activityRepo:  activityRepo,
	}
}

func (b *BusinessServer) GetHome(ctx context.Context, req *pb.GetHomeRequest) (*pb.GetHomeResponse, error) {
	log.Infof("GetHome for network %q, window %+v", req.GetNetworkId(), req.GetWindow())

	now := time.Now().UTC()

	win, err := resolveWindow(req.GetWindow(), now)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid window: %v", err)
	}

	revenue, err := b.salesRepo.RevenueBetween(req.GetNetworkId(), win.From, win.To)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "revenue")
	}

	prevRevenue, err := b.salesRepo.RevenueBetween(req.GetNetworkId(), win.PrevFrom, win.PrevTo)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "revenue")
	}

	packages, _, err := b.packageRepo.ListPackages(0, 0)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "packages")
	}

	active := activeSubscribers(packages)

	uptime, err := b.siteRepo.SiteUptime("", win.From, win.To)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "uptime")
	}

	siteRollups, err := b.siteRepo.SiteRollups("", win.From, win.To)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "site rollups")
	}

	dataSold := 0.0
	for _, r := range siteRollups {
		dataSold += r.DataUsedMb
	}

	kpis := []*pb.Kpi{
		makeKpi("revenue_today", revenue, formatMoney(revenue),
			pctDelta(revenue, prevRevenue), win.Period, now),
		makeKpi("active_customers", float64(active), formatCount(active), 0, win.Period, now),
		makeKpi("data_sold", dataSold, fmt.Sprintf("%.1f MB", dataSold), 0, win.Period, now),
		makeKpi("network_uptime", uptime, formatPercent(uptime), 0, win.Period, now),
	}

	sites, _, err := b.siteRepo.ListSites(req.GetNetworkId(), 0, 0)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sites")
	}

	revenueBySite, err := b.salesRepo.RevenueBySite(req.GetNetworkId(), win.From, win.To)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "revenue by site")
	}

	siteRevenue := map[string]float64{}
	for _, r := range revenueBySite {
		siteRevenue[r.Id] = r.Value
	}

	siteCustomers := map[string]uint32{}
	for _, r := range siteRollups {
		siteCustomers[r.SiteId.String()] = r.Customers
	}

	siteSummaries := make([]*pb.SiteSummary, 0, len(sites))
	for _, s := range sites {
		id := s.SiteId.String()
		siteSummaries = append(siteSummaries, &pb.SiteSummary{
			SiteId:    id,
			Name:      s.Name,
			Status:    s.Status,
			Revenue:   siteRevenue[id],
			Customers: siteCustomers[id],
		})
	}

	topPackages, err := b.salesRepo.RevenueByPackage(req.GetNetworkId(), win.From, win.To)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "revenue by package")
	}

	topPkgValues := make([]*pb.NamedValue, 0, len(topPackages))
	for _, p := range topPackages {
		topPkgValues = append(topPkgValues, &pb.NamedValue{
			Name:  p.Name,
			Id:    p.Id,
			Value: p.Value,
		})
	}

	events, err := b.activityRepo.Recent(10)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "activity")
	}

	activity := make([]*pb.ActivityItem, 0, len(events))
	for _, e := range events {
		activity = append(activity, &pb.ActivityItem{
			RoutingKey:  e.RoutingKey,
			Description: e.RoutingKey,
			OccurredAt:  timestamppb.New(e.OccurredAt),
		})
	}

	return &pb.GetHomeResponse{
		Kpis:           kpis,
		Sites:          siteSummaries,
		TopPackages:    topPkgValues,
		RecentActivity: activity,
	}, nil
}

func (b *BusinessServer) GetSalesOverview(ctx context.Context, req *pb.GetSalesOverviewRequest) (*pb.GetSalesOverviewResponse, error) {
	log.Infof("GetSalesOverview for network %q, window %+v", req.GetNetworkId(), req.GetWindow())

	now := time.Now().UTC()

	win, err := resolveWindow(req.GetWindow(), now)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid window: %v", err)
	}

	networkId := req.GetNetworkId()

	revenue, err := b.salesRepo.RevenueBetween(networkId, win.From, win.To)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "revenue")
	}

	prevRevenue, err := b.salesRepo.RevenueBetween(networkId, win.PrevFrom, win.PrevTo)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "revenue")
	}

	purchases, err := b.salesRepo.PurchasesBetween(networkId, win.From, win.To)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "purchases")
	}

	prevPurchases, err := b.salesRepo.PurchasesBetween(networkId, win.PrevFrom, win.PrevTo)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "purchases")
	}

	paidCustomers, err := b.salesRepo.PaidCustomersBetween(networkId, win.From, win.To)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "paid customers")
	}

	prevPaidCustomers, err := b.salesRepo.PaidCustomersBetween(networkId, win.PrevFrom, win.PrevTo)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "paid customers")
	}

	avg := avgPurchase(revenue, purchases)
	prevAvg := avgPurchase(prevRevenue, prevPurchases)

	kpis := []*pb.Kpi{
		makeKpi("revenue", revenue, formatMoney(revenue),
			pctDelta(revenue, prevRevenue), win.Period, now),
		makeKpi("purchases", float64(purchases), formatCount(purchases),
			pctDelta(float64(purchases), float64(prevPurchases)), win.Period, now),
		makeKpi("avg_purchase", avg, formatMoney(avg),
			pctDelta(avg, prevAvg), win.Period, now),
		makeKpi("paid_customers", float64(paidCustomers), formatCount(paidCustomers),
			pctDelta(float64(paidCustomers), float64(prevPaidCustomers)), win.Period, now),
	}

	trend, err := b.salesRepo.RevenueTrendDaily(networkId, win.From, win.To)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "revenue trend")
	}

	points := make([]*pb.Point, 0, len(trend))
	for _, d := range trend {
		points = append(points, &pb.Point{
			Time:  timestamppb.New(d.Day),
			Value: d.Value,
		})
	}

	bySite, err := b.salesRepo.RevenueBySite(networkId, win.From, win.To)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "revenue by site")
	}

	byPackage, err := b.salesRepo.RevenueByPackage(networkId, win.From, win.To)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "revenue by package")
	}

	return &pb.GetSalesOverviewResponse{
		Kpis: kpis,
		RevenueTrend: &pb.TimeSeries{
			Key:    "revenue",
			Points: points,
		},
		RevenueBySite:    toNamedValues(bySite),
		RevenueByPackage: toNamedValues(byPackage),
	}, nil
}

func (b *BusinessServer) GetPackagePerformance(ctx context.Context, req *pb.GetPackagePerformanceRequest) (*pb.GetPackagePerformanceResponse, error) {
	log.Infof("GetPackagePerformance for network %q, window %+v", req.GetNetworkId(), req.GetWindow())

	now := time.Now().UTC()

	win, err := resolveWindow(req.GetWindow(), now)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid window: %v", err)
	}

	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	if pageSize == 0 {
		pageSize = defaultPageSize
	}

	packages, total, err := b.packageRepo.ListPackages(page, pageSize)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "packages")
	}

	// MRR/ARPU use the full package population, not the current page.
	allPackages, _, err := b.packageRepo.ListPackages(0, 0)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "packages")
	}

	rollups, err := b.packageRepo.PackageRollups(win.From, win.To)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "package rollups")
	}

	revenue, err := b.salesRepo.RevenueBetween(req.GetNetworkId(), win.From, win.To)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "revenue")
	}

	mrrValue := mrr(allPackages)
	active := activeSubscribers(allPackages)
	arpuValue := arpu(revenue, active)

	_, topRevenue, topShare, _ := topPlanByRevenue(rollups)

	kpis := []*pb.Kpi{
		makeKpi("mrr", mrrValue, formatMoney(mrrValue), 0, win.Period, now),
		makeKpi("arpu", arpuValue, formatMoney(arpuValue), 0, win.Period, now),
		makeKpi("top_plan_revenue", topRevenue, formatMoney(topRevenue), 0, win.Period, now),
		makeKpi("top_plan_share", topShare, formatPercent(topShare), 0, win.Period, now),
	}

	// Aggregate per-package rollups for the window.
	type pkgAgg struct {
		soldCount  uint32
		revenue    float64
		dataUsedMb float64
	}

	aggs := map[string]pkgAgg{}
	for _, r := range rollups {
		id := r.PackageId.String()
		a := aggs[id]
		a.soldCount += r.SoldCount
		a.revenue += r.Revenue
		a.dataUsedMb += r.DataUsedMb
		aggs[id] = a
	}

	rows := make([]*pb.PackageRow, 0, len(packages))
	mix := make([]*pb.NamedValue, 0, len(packages))
	for _, p := range packages {
		id := p.PackageId.String()
		a := aggs[id]

		rows = append(rows, &pb.PackageRow{
			PackageId:         id,
			Name:              p.Name,
			Price:             p.Price,
			Validity:          fmt.Sprintf("%d days", p.DurationDays),
			DataQuota:         fmt.Sprintf("%.0f MB", p.DataQuotaMb),
			Status:            p.Status,
			SoldCount:         a.soldCount,
			Revenue:           a.revenue,
			DataUsed:          a.dataUsedMb,
			ActiveSubscribers: p.ActiveSubscribers,
		})

		mix = append(mix, &pb.NamedValue{
			Name:  p.Name,
			Id:    id,
			Value: a.revenue,
		})
	}

	return &pb.GetPackagePerformanceResponse{
		Kpis:       kpis,
		Packages:   rows,
		RevenueMix: mix,
		Meta:       makeMeta(uint32(total), req.GetPage(), uint32(pageSize)),
	}, nil
}

func (b *BusinessServer) GetBillingSummary(ctx context.Context, req *pb.GetBillingSummaryRequest) (*pb.GetBillingSummaryResponse, error) {
	log.Infof("GetBillingSummary, window %+v", req.GetWindow())

	now := time.Now().UTC()

	win, err := resolveWindow(req.GetWindow(), now)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid window: %v", err)
	}

	snap, err := b.billingRepo.GetBillingSnapshot()
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "billing snapshot")
	}

	rollups, err := b.billingRepo.InvoiceRollups(win.From, win.To)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "invoice rollups")
	}

	var lastInvoiceAmount float64
	if len(rollups) > 0 {
		// InvoiceRollups orders by day DESC; the first row is the latest.
		lastInvoiceAmount = rollups[0].InvoicedAmount
	}

	kpis := []*pb.Kpi{
		makeKpi("current_balance", snap.Balance, formatMoney(snap.Balance), 0, win.Period, now),
		makeKpi("last_invoice_amount", lastInvoiceAmount, formatMoney(lastInvoiceAmount), 0, win.Period, now),
	}

	invoices := make([]*pb.InvoiceRow, 0, len(rollups))
	for _, r := range rollups {
		invoices = append(invoices, &pb.InvoiceRow{
			InvoiceId:   r.Day.Format("2006-01-02"),
			Amount:      r.InvoicedAmount,
			Status:      "generated",
			GeneratedAt: timestamppb.New(r.Day),
		})
	}

	var lastInvoiceDate *timestamppb.Timestamp
	if snap.LastInvoiceAt != nil {
		lastInvoiceDate = timestamppb.New(*snap.LastInvoiceAt)
	}

	pageSize := req.GetPageSize()
	if pageSize == 0 {
		pageSize = defaultPageSize
	}

	return &pb.GetBillingSummaryResponse{
		Kpis:            kpis,
		Invoices:        invoices,
		LastInvoiceDate: lastInvoiceDate,
		Meta:            makeMeta(uint32(len(invoices)), req.GetPage(), pageSize),
	}, nil
}

func (b *BusinessServer) GetSites(ctx context.Context, req *pb.GetSitesRequest) (*pb.GetSitesResponse, error) {
	log.Infof("GetSites for network %q, window %+v", req.GetNetworkId(), req.GetWindow())

	now := time.Now().UTC()

	win, err := resolveWindow(req.GetWindow(), now)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid window: %v", err)
	}

	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	if pageSize == 0 {
		pageSize = defaultPageSize
	}

	sites, total, err := b.siteRepo.ListSites(req.GetNetworkId(), page, pageSize)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sites")
	}

	rows := make([]*pb.BusinessSiteRow, 0, len(sites))
	for _, s := range sites {
		row, err := b.buildSiteRow(&s, win, now)
		if err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}

	return &pb.GetSitesResponse{
		Sites: rows,
		Meta:  makeMeta(uint32(total), req.GetPage(), uint32(pageSize)),
	}, nil
}

func (b *BusinessServer) GetSite(ctx context.Context, req *pb.GetSiteRequest) (*pb.GetSiteResponse, error) {
	log.Infof("GetSite %q, window %+v", req.GetSiteId(), req.GetWindow())

	if req.GetSiteId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "site_id is required")
	}

	now := time.Now().UTC()

	win, err := resolveWindow(req.GetWindow(), now)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid window: %v", err)
	}

	site, err := b.siteRepo.GetSite(req.GetSiteId())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "site")
	}

	row, err := b.buildSiteRow(site, win, now)
	if err != nil {
		return nil, err
	}

	rollups, err := b.siteRepo.SiteRollups(req.GetSiteId(), win.From, win.To)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "site rollups")
	}

	points := make([]*pb.Point, 0, len(rollups))
	for _, r := range rollups {
		points = append(points, &pb.Point{
			Time:  timestamppb.New(r.Day),
			Value: r.Revenue,
		})
	}

	kpis := []*pb.Kpi{
		makeKpi("site_revenue", row.Revenue, formatMoney(row.Revenue), 0, win.Period, now),
		makeKpi("site_customers", float64(row.Customers), formatCount(row.Customers), 0, win.Period, now),
		makeKpi("site_uptime", row.Uptime, formatPercent(row.Uptime), 0, win.Period, now),
	}

	return &pb.GetSiteResponse{
		Site: row,
		Kpis: kpis,
		RevenueTrend: &pb.TimeSeries{
			Key:    "site_revenue",
			Points: points,
		},
	}, nil
}

func (b *BusinessServer) GetInventoryReadiness(ctx context.Context, req *pb.GetInventoryReadinessRequest) (*pb.GetInventoryReadinessResponse, error) {
	log.Infof("GetInventoryReadiness for network %q", req.GetNetworkId())

	now := time.Now().UTC()

	availableSims, activeSims, err := b.inventoryRepo.SimCounts()
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sim counts")
	}

	availableNodes, deployedNodes, err := b.inventoryRepo.NodeCounts()
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "node counts")
	}

	kpis := []*pb.Kpi{
		makeKpi("available_sims", float64(availableSims), formatCount(availableSims), 0, "", now),
		makeKpi("active_sims", float64(activeSims), formatCount(activeSims), 0, "", now),
		makeKpi("available_nodes", float64(availableNodes), formatCount(availableNodes), 0, "", now),
		makeKpi("deployed_nodes", float64(deployedNodes), formatCount(deployedNodes), 0, "", now),
	}

	return &pb.GetInventoryReadinessResponse{
		Kpis: kpis,
	}, nil
}

// buildSiteRow assembles a BusinessSiteRow from the site snapshot plus
// windowed rollups (revenue, customers, data, uptime).
func (b *BusinessServer) buildSiteRow(site *db.SiteSnapshot, win *ResolvedWindow, now time.Time) (*pb.BusinessSiteRow, error) {
	id := site.SiteId.String()

	rollups, err := b.siteRepo.SiteRollups(id, win.From, win.To)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "site rollups")
	}

	var revenue, dataUsed float64
	var customers uint32
	for _, r := range rollups {
		revenue += r.Revenue
		dataUsed += r.DataUsedMb
		if r.Customers > customers {
			customers = r.Customers
		}
	}

	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, win.Location)

	todayRollups, err := b.siteRepo.SiteRollups(id, startOfDay, now)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "site rollups")
	}

	var revenueToday float64
	for _, r := range todayRollups {
		revenueToday += r.Revenue
	}

	uptime, err := b.siteRepo.SiteUptime(id, win.From, win.To)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "site uptime")
	}

	issue := ""
	if site.Status != "online" {
		issue = site.Status
	}

	return &pb.BusinessSiteRow{
		SiteId:       id,
		Name:         site.Name,
		Status:       site.Status,
		Revenue:      revenue,
		RevenueToday: revenueToday,
		Customers:    customers,
		DataUsed:     dataUsed,
		Uptime:       uptime,
		TopPackage:   "",
		Issue:        issue,
		Latitude:     site.Latitude,
		Longitude:    site.Longitude,
	}, nil
}

func toNamedValues(rows []db.NamedAmount) []*pb.NamedValue {
	out := make([]*pb.NamedValue, 0, len(rows))
	for _, r := range rows {
		out = append(out, &pb.NamedValue{
			Name:  r.Name,
			Id:    r.Id,
			Value: r.Value,
		})
	}

	return out
}

func makeMeta(count, page, size uint32) *pb.Meta {
	if page == 0 {
		page = 1
	}

	var pages uint32
	if size > 0 {
		pages = (count + size - 1) / size
	}

	return &pb.Meta{
		Count: count,
		Page:  page,
		Size:  size,
		Pages: pages,
	}
}
