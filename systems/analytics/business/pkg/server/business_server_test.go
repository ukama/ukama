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
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"

	"github.com/ukama/ukama/systems/analytics/business/mocks"
	pb "github.com/ukama/ukama/systems/analytics/business/pb/gen"
	"github.com/ukama/ukama/systems/analytics/business/pkg/db"
	"github.com/ukama/ukama/systems/common/uuid"
)

type bizMocks struct {
	sales     *mocks.SalesRepo
	pkg       *mocks.PackageRepo
	site      *mocks.SiteRepo
	billing   *mocks.BillingRepo
	inventory *mocks.InventoryRepo
	activity  *mocks.ActivityRepo
}

func newBizMocks() *bizMocks {
	return &bizMocks{
		sales:     &mocks.SalesRepo{},
		pkg:       &mocks.PackageRepo{},
		site:      &mocks.SiteRepo{},
		billing:   &mocks.BillingRepo{},
		inventory: &mocks.InventoryRepo{},
		activity:  &mocks.ActivityRepo{},
	}
}

func (m *bizMocks) server() *BusinessServer {
	return NewBusinessServer("test-org", m.sales, m.pkg, m.site, m.billing, m.inventory, m.activity)
}

func sampleSite() db.SiteSnapshot {
	return db.SiteSnapshot{SiteId: uuid.NewV4(), Name: "Site One", Status: "online", Latitude: 1, Longitude: 2}
}

func samplePackage() db.PackageSnapshot {
	return db.PackageSnapshot{PackageId: uuid.NewV4(), Name: "Starter", Price: 5, DurationDays: 30, DataQuotaMb: 1024, Status: "active", ActiveSubscribers: 12}
}

func TestGetHome(t *testing.T) {
	m := newBizMocks()
	site := sampleSite()

	m.sales.On("RevenueBetween", mock.Anything, mock.Anything, mock.Anything).Return(100.0, nil)
	m.pkg.On("ListPackages", 0, 0).Return([]db.PackageSnapshot{samplePackage()}, int64(1), nil)
	m.site.On("SiteUptime", mock.Anything, mock.Anything, mock.Anything).Return(99.5, nil)
	m.site.On("SiteRollups", mock.Anything, mock.Anything, mock.Anything).
		Return([]db.BusinessSiteRollupDaily{{Day: time.Now(), SiteId: site.SiteId, Revenue: 50, DataUsedMb: 10, Customers: 3}}, nil)
	m.site.On("ListSites", mock.Anything, 0, 0).Return([]db.SiteSnapshot{site}, int64(1), nil)
	m.sales.On("RevenueBySite", mock.Anything, mock.Anything, mock.Anything).
		Return([]db.NamedAmount{{Id: site.SiteId.String(), Name: "Site One", Value: 50}}, nil)
	m.sales.On("RevenueByPackage", mock.Anything, mock.Anything, mock.Anything).
		Return([]db.NamedAmount{{Id: "pkg-1", Name: "Starter", Value: 30}}, nil)
	m.activity.On("Recent", 10).Return([]db.EventLog{{RoutingKey: "event.x", OccurredAt: time.Now()}}, nil)

	resp, err := m.server().GetHome(context.Background(), &pb.GetHomeRequest{Window: &pb.Window{Period: PeriodToday}})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Kpis, 4)
	assert.Len(t, resp.Sites, 1)
}

func TestGetHome_InvalidWindow(t *testing.T) {
	m := newBizMocks()
	_, err := m.server().GetHome(context.Background(), &pb.GetHomeRequest{Window: &pb.Window{Period: PeriodCustom}})
	assert.Error(t, err)
}

func TestGetHome_RepoError(t *testing.T) {
	m := newBizMocks()
	m.sales.On("RevenueBetween", mock.Anything, mock.Anything, mock.Anything).Return(0.0, errors.New("db down"))

	_, err := m.server().GetHome(context.Background(), &pb.GetHomeRequest{Window: &pb.Window{Period: PeriodToday}})
	assert.Error(t, err)
}

func TestGetPackagePerformance(t *testing.T) {
	m := newBizMocks()
	pkg := samplePackage()

	m.pkg.On("ListPackages", mock.Anything, mock.Anything).Return([]db.PackageSnapshot{pkg}, int64(1), nil)
	m.pkg.On("PackageRollups", mock.Anything, mock.Anything).
		Return([]db.BusinessPackageRollupDaily{{Day: time.Now(), PackageId: pkg.PackageId, SoldCount: 5, Revenue: 25, DataUsedMb: 7}}, nil)
	m.sales.On("RevenueBetween", mock.Anything, mock.Anything, mock.Anything).Return(100.0, nil)

	resp, err := m.server().GetPackagePerformance(context.Background(), &pb.GetPackagePerformanceRequest{
		Window: &pb.Window{Period: PeriodWeek}, Page: 1, PageSize: 10,
	})

	assert.NoError(t, err)
	assert.Len(t, resp.Packages, 1)
	assert.Len(t, resp.Kpis, 4)
}

func TestGetBillingSummary(t *testing.T) {
	m := newBizMocks()
	last := time.Now()

	m.billing.On("GetBillingSnapshot").Return(&db.BillingSnapshot{Balance: 250, LastInvoiceAt: &last}, nil)
	m.billing.On("InvoiceRollups", mock.Anything, mock.Anything).
		Return([]db.BusinessBillingRollupDaily{{Day: time.Now(), InvoicedAmount: 80}}, nil)

	resp, err := m.server().GetBillingSummary(context.Background(), &pb.GetBillingSummaryRequest{Window: &pb.Window{Period: PeriodMonth}})

	assert.NoError(t, err)
	assert.Len(t, resp.Invoices, 1)
	assert.NotNil(t, resp.LastInvoiceDate)
}

func TestGetSites(t *testing.T) {
	m := newBizMocks()
	site := sampleSite()

	m.site.On("ListSites", mock.Anything, mock.Anything, mock.Anything).Return([]db.SiteSnapshot{site}, int64(1), nil)
	m.site.On("SiteRollups", mock.Anything, mock.Anything, mock.Anything).
		Return([]db.BusinessSiteRollupDaily{{Day: time.Now(), SiteId: site.SiteId, Revenue: 50, Customers: 3, DataUsedMb: 10}}, nil)
	m.site.On("SiteUptime", mock.Anything, mock.Anything, mock.Anything).Return(98.0, nil)

	resp, err := m.server().GetSites(context.Background(), &pb.GetSitesRequest{Window: &pb.Window{Period: PeriodWeek}, Page: 1, PageSize: 10})

	assert.NoError(t, err)
	assert.Len(t, resp.Sites, 1)
}

func TestGetSite(t *testing.T) {
	m := newBizMocks()
	site := sampleSite()

	m.site.On("GetSite", mock.Anything).Return(&site, nil)
	m.site.On("SiteRollups", mock.Anything, mock.Anything, mock.Anything).
		Return([]db.BusinessSiteRollupDaily{{Day: time.Now(), SiteId: site.SiteId, Revenue: 50, Customers: 3}}, nil)
	m.site.On("SiteUptime", mock.Anything, mock.Anything, mock.Anything).Return(97.0, nil)

	resp, err := m.server().GetSite(context.Background(), &pb.GetSiteRequest{SiteId: site.SiteId.String(), Window: &pb.Window{Period: PeriodWeek}})

	assert.NoError(t, err)
	assert.NotNil(t, resp.Site)
	assert.Len(t, resp.Kpis, 3)
}

func TestGetSite_MissingId(t *testing.T) {
	m := newBizMocks()
	_, err := m.server().GetSite(context.Background(), &pb.GetSiteRequest{SiteId: "", Window: &pb.Window{Period: PeriodWeek}})
	assert.Error(t, err)
}

func TestGetInventoryReadiness(t *testing.T) {
	m := newBizMocks()

	m.inventory.On("SimCounts").Return(uint32(100), uint32(40), nil)
	m.inventory.On("NodeCounts").Return(uint32(20), uint32(8), nil)

	resp, err := m.server().GetInventoryReadiness(context.Background(), &pb.GetInventoryReadinessRequest{})

	assert.NoError(t, err)
	assert.Len(t, resp.Kpis, 4)
}

func TestGetInventoryReadiness_Error(t *testing.T) {
	m := newBizMocks()
	m.inventory.On("SimCounts").Return(uint32(0), uint32(0), errors.New("db down"))

	_, err := m.server().GetInventoryReadiness(context.Background(), &pb.GetInventoryReadinessRequest{})
	assert.Error(t, err)
}
