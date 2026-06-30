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

	"github.com/ukama/ukama/systems/analytics/customer/mocks"
	pb "github.com/ukama/ukama/systems/analytics/customer/pb/gen"
	"github.com/ukama/ukama/systems/analytics/customer/pkg/db"
	"github.com/ukama/ukama/systems/common/uuid"
)

type custMocks struct {
	customer *mocks.CustomerRepo
	sim      *mocks.SimRepo
	support  *mocks.SupportRepo
}

func newCustMocks() *custMocks {
	return &custMocks{
		customer: &mocks.CustomerRepo{},
		sim:      &mocks.SimRepo{},
		support:  &mocks.SupportRepo{},
	}
}

func (m *custMocks) server() *CustomerServer {
	return NewCustomerServer("test-org", m.customer, m.sim, m.support, 50)
}

func sampleCustomer() db.CustomerSnapshot {
	return db.CustomerSnapshot{
		CustomerId: uuid.NewV4(),
		Name:       "Jane",
		Email:      "jane@x.io",
		Status:     "active",
		PackageId:  uuid.NewV4(),
		SimStatus:  "active",
		SiteId:     uuid.Nil,
	}
}

func TestCustomer_GetOverview(t *testing.T) {
	m := newCustMocks()
	m.customer.On("Counts", mock.Anything, mock.Anything, mock.Anything).
		Return(uint32(100), uint32(80), uint32(10), uint32(5), uint32(2), nil)

	resp, err := m.server().GetOverview(context.Background(), &pb.GetOverviewRequest{Window: &pb.Window{Period: PeriodWeek}})

	assert.NoError(t, err)
	assert.Len(t, resp.Kpis, 5)
}

func TestCustomer_GetOverview_BadNetwork(t *testing.T) {
	m := newCustMocks()
	_, err := m.server().GetOverview(context.Background(), &pb.GetOverviewRequest{NetworkId: "not-a-uuid"})
	assert.Error(t, err)
}

func TestCustomer_GetOverview_RepoError(t *testing.T) {
	m := newCustMocks()
	m.customer.On("Counts", mock.Anything, mock.Anything, mock.Anything).
		Return(uint32(0), uint32(0), uint32(0), uint32(0), uint32(0), errors.New("db down"))

	_, err := m.server().GetOverview(context.Background(), &pb.GetOverviewRequest{})
	assert.Error(t, err)
}

func TestCustomer_List(t *testing.T) {
	m := newCustMocks()
	m.customer.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return([]db.CustomerSnapshot{sampleCustomer()}, int64(1), nil)
	m.customer.On("SiteNames", mock.Anything).Return(map[uuid.UUID]string{}, nil)

	resp, err := m.server().List(context.Background(), &pb.ListRequest{Page: 1, PageSize: 10})

	assert.NoError(t, err)
	assert.Len(t, resp.Customers, 1)
}

func TestCustomer_Search_EmptyQuery(t *testing.T) {
	m := newCustMocks()
	_, err := m.server().Search(context.Background(), &pb.SearchRequest{Query: ""})
	assert.Error(t, err)
}

func TestCustomer_Search(t *testing.T) {
	m := newCustMocks()
	m.customer.On("Search", "jane", mock.Anything, mock.Anything, mock.Anything).
		Return([]db.CustomerSnapshot{sampleCustomer()}, int64(1), nil)
	m.customer.On("SiteNames", mock.Anything).Return(map[uuid.UUID]string{}, nil)

	resp, err := m.server().Search(context.Background(), &pb.SearchRequest{Query: "jane"})

	assert.NoError(t, err)
	assert.Len(t, resp.Customers, 1)
}

func TestCustomer_Get(t *testing.T) {
	m := newCustMocks()
	cust := sampleCustomer()
	m.customer.On("Get", mock.Anything).Return(&cust, nil)
	m.customer.On("SiteNames", mock.Anything).Return(map[uuid.UUID]string{}, nil)
	m.customer.On("UsageBetween", mock.Anything, mock.Anything, mock.Anything).Return(123.4, nil)
	m.customer.On("PackageIntervals", mock.Anything).
		Return([]db.CustomerPackageInterval{{CustomerId: cust.CustomerId, PackageId: cust.PackageId, State: "active", StartAt: time.Now()}}, nil)

	resp, err := m.server().Get(context.Background(), &pb.GetRequest{CustomerId: cust.CustomerId.String(), Window: &pb.Window{Period: PeriodWeek}})

	assert.NoError(t, err)
	assert.NotNil(t, resp.Customer)
	assert.Len(t, resp.Kpis, 2)
	assert.Len(t, resp.PackageHistory, 1)
}

func TestCustomer_Get_BadId(t *testing.T) {
	m := newCustMocks()
	_, err := m.server().Get(context.Background(), &pb.GetRequest{CustomerId: "bad"})
	assert.Error(t, err)
}

func TestCustomer_GetSupport(t *testing.T) {
	m := newCustMocks()
	cust := sampleCustomer()
	m.customer.On("Get", mock.Anything).Return(&cust, nil)
	m.customer.On("SiteNames", mock.Anything).Return(map[uuid.UUID]string{}, nil)
	m.customer.On("UsageBetween", mock.Anything, mock.Anything, mock.Anything).Return(10.0, nil)
	m.support.On("RecentActivityFor", mock.Anything, mock.Anything).
		Return([]db.EventLog{{RoutingKey: "event.x", OccurredAt: time.Now()}}, nil)

	resp, err := m.server().GetSupport(context.Background(), &pb.GetSupportRequest{CustomerId: cust.CustomerId.String()})

	assert.NoError(t, err)
	assert.NotNil(t, resp.Customer)
}

func TestCustomer_GetSims(t *testing.T) {
	m := newCustMocks()
	allocated := time.Now()
	m.sim.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return([]db.SimSnapshot{{SimId: "sim-1", Iccid: "8910", Status: "active", CustomerId: uuid.NewV4(), BatchId: "b1", AllocatedAt: &allocated}}, int64(1), nil)

	resp, err := m.server().GetSims(context.Background(), &pb.GetSimsRequest{Page: 1, PageSize: 10})

	assert.NoError(t, err)
	assert.Len(t, resp.Sims, 1)
}

func TestCustomer_GetSimPool(t *testing.T) {
	m := newCustMocks()
	uploaded := time.Now()
	m.sim.On("PoolCounts").Return(uint32(200), uint32(10), uint32(150), uint32(30), uint32(5), uint32(5), nil)
	m.sim.On("Batches").Return([]db.SimBatchSnapshot{{BatchId: "b1", Quantity: 100, Assigned: 60, UploadedAt: &uploaded}}, nil)

	resp, err := m.server().GetSimPool(context.Background(), &pb.GetSimPoolRequest{})

	assert.NoError(t, err)
	assert.Len(t, resp.Batches, 1)
	// available (10) < threshold (50) => low_stock kpi present
	assert.NotEmpty(t, resp.Kpis)
}

func TestCustomer_GetSimPool_RepoError(t *testing.T) {
	m := newCustMocks()
	m.sim.On("PoolCounts").Return(uint32(0), uint32(0), uint32(0), uint32(0), uint32(0), uint32(0), errors.New("db down"))

	_, err := m.server().GetSimPool(context.Background(), &pb.GetSimPoolRequest{})
	assert.Error(t, err)
}
