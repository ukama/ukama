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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	pb "github.com/ukama/ukama/systems/analytics/business/pb/gen"
	"github.com/ukama/ukama/systems/analytics/business/pkg/db"
)

// stubSalesRepo is a hand-written stub (no mockery) returning canned values
// keyed by window: the "current" window is any interval ending after pivot,
// the "previous" window any interval ending at or before pivot.
type stubSalesRepo struct {
	pivot         time.Time
	revenue       float64
	prevRevenue   float64
	purchases     uint32
	prevPurchases uint32
	paid          uint32
	prevPaid      uint32
}

func (s *stubSalesRepo) isCurrent(to time.Time) bool {
	return to.After(s.pivot)
}

func (s *stubSalesRepo) RevenueBetween(networkId string, from, to time.Time) (float64, error) {
	if s.isCurrent(to) {
		return s.revenue, nil
	}

	return s.prevRevenue, nil
}

func (s *stubSalesRepo) PurchasesBetween(networkId string, from, to time.Time) (uint32, error) {
	if s.isCurrent(to) {
		return s.purchases, nil
	}

	return s.prevPurchases, nil
}

func (s *stubSalesRepo) PaidCustomersBetween(networkId string, from, to time.Time) (uint32, error) {
	if s.isCurrent(to) {
		return s.paid, nil
	}

	return s.prevPaid, nil
}

func (s *stubSalesRepo) RevenueTrendDaily(networkId string, from, to time.Time) ([]db.DayValue, error) {
	return []db.DayValue{{Day: from, Value: s.revenue}}, nil
}

func (s *stubSalesRepo) RevenueBySite(networkId string, from, to time.Time) ([]db.NamedAmount, error) {
	return []db.NamedAmount{{Id: "site-1", Name: "Site One", Value: s.revenue}}, nil
}

func (s *stubSalesRepo) RevenueByPackage(networkId string, from, to time.Time) ([]db.NamedAmount, error) {
	return []db.NamedAmount{{Id: "pkg-1", Name: "Starter", Value: s.revenue}}, nil
}

func kpiByKey(kpis []*pb.Kpi, key string) *pb.Kpi {
	for _, k := range kpis {
		if k.Key == key {
			return k
		}
	}

	return nil
}

func TestGetSalesOverview_KpisAndDeltas(t *testing.T) {
	startOfToday := time.Now().UTC().Truncate(24 * time.Hour)

	sales := &stubSalesRepo{
		pivot:         startOfToday,
		revenue:       200,
		prevRevenue:   100,
		purchases:     8,
		prevPurchases: 10,
		paid:          4,
		prevPaid:      4,
	}

	s := NewBusinessServer("test-org", sales, nil, nil, nil, nil, nil)

	resp, err := s.GetSalesOverview(context.TODO(), &pb.GetSalesOverviewRequest{
		Window: &pb.Window{Period: PeriodToday},
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	revenue := kpiByKey(resp.Kpis, "revenue")
	assert.NotNil(t, revenue)
	assert.Equal(t, 200.0, revenue.Value)
	assert.Equal(t, "$200.00", revenue.Formatted)
	assert.InDelta(t, 100.0, revenue.Delta, 1e-9, "revenue doubled vs prev window")

	purchases := kpiByKey(resp.Kpis, "purchases")
	assert.NotNil(t, purchases)
	assert.Equal(t, 8.0, purchases.Value)
	assert.InDelta(t, -20.0, purchases.Delta, 1e-9)

	avg := kpiByKey(resp.Kpis, "avg_purchase")
	assert.NotNil(t, avg)
	assert.InDelta(t, 25.0, avg.Value, 1e-9) // 200/8
	// prev avg = 100/10 = 10; delta = (25-10)/10 = 150%
	assert.InDelta(t, 150.0, avg.Delta, 1e-9)

	paid := kpiByKey(resp.Kpis, "paid_customers")
	assert.NotNil(t, paid)
	assert.Equal(t, 4.0, paid.Value)
	assert.InDelta(t, 0.0, paid.Delta, 1e-9, "no change vs prev window")

	assert.NotNil(t, resp.RevenueTrend)
	assert.Len(t, resp.RevenueTrend.Points, 1)
	assert.Len(t, resp.RevenueBySite, 1)
	assert.Len(t, resp.RevenueByPackage, 1)
}

func TestGetSalesOverview_InvalidWindow(t *testing.T) {
	s := NewBusinessServer("test-org", &stubSalesRepo{}, nil, nil, nil, nil, nil)

	_, err := s.GetSalesOverview(context.TODO(), &pb.GetSalesOverviewRequest{
		Window: &pb.Window{Period: "fortnight"},
	})
	assert.Error(t, err)
}
