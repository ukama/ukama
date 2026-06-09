/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package server

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/ukama/ukama/systems/analytics/network/pb/gen"
	"github.com/ukama/ukama/systems/analytics/network/pkg/db"
)

func kpi(key string, value float64, formatted string) *pb.Kpi {
	return &pb.Kpi{
		Key:       key,
		Value:     value,
		Formatted: formatted,
		AsOf:      timestamppb.New(time.Now().UTC()),
	}
}

func countKpi(key string, value int64) *pb.Kpi {
	return kpi(key, float64(value), fmt.Sprintf("%d", value))
}

func percentKpi(key string, value float64) *pb.Kpi {
	return kpi(key, value, fmt.Sprintf("%.1f%%", value))
}

func ratioKpi(key string, num, den int64) *pb.Kpi {
	return kpi(key, float64(num), fmt.Sprintf("%d/%d", num, den))
}

func moneyKpi(key string, value float64) *pb.Kpi {
	return kpi(key, value, fmt.Sprintf("%.2f", value))
}

func floatKpi(key string, value float64, unit string) *pb.Kpi {
	return kpi(key, value, fmt.Sprintf("%.2f %s", value, unit))
}

func pbMeta(count int64, page, pageSize uint32) *pb.Meta {
	if page < 1 {
		page = 1
	}

	pages := uint32(1)
	if pageSize > 0 {
		pages = uint32((count + int64(pageSize) - 1) / int64(pageSize))
		if pages == 0 {
			pages = 1
		}
	}

	return &pb.Meta{
		Count: uint32(count),
		Page:  page,
		Size:  pageSize,
		Pages: pages,
	}
}

func pbAlarm(a *db.AlarmEvent) *pb.AlarmRow {
	row := &pb.AlarmRow{
		AlarmId:           a.AlarmId,
		Severity:          a.Severity,
		State:             a.State,
		ResourceType:      a.ResourceType,
		ResourceId:        a.ResourceId,
		Description:       a.Description,
		CustomersAffected: a.CustomersAffected,
		RevenueAtRisk:     a.RevenueAtRisk,
		RecommendedAction: a.RecommendedAction,
		OpenedAt:          timestamppb.New(a.OpenedAt),
	}

	if a.ClosedAt != nil {
		row.ClosedAt = timestamppb.New(*a.ClosedAt)
	}

	return row
}

func pbAlarms(alarms []db.AlarmEvent) []*pb.AlarmRow {
	rows := make([]*pb.AlarmRow, 0, len(alarms))
	for i := range alarms {
		rows = append(rows, pbAlarm(&alarms[i]))
	}

	return rows
}

func pbEvent(e *db.EventLog) *pb.EventRow {
	return &pb.EventRow{
		RoutingKey:  e.RoutingKey,
		Description: string(e.Payload),
		OccurredAt:  timestamppb.New(e.OccurredAt),
	}
}

func pbEvents(events []db.EventLog) []*pb.EventRow {
	rows := make([]*pb.EventRow, 0, len(events))
	for i := range events {
		rows = append(rows, pbEvent(&events[i]))
	}

	return rows
}

func point(t time.Time, v float64) *pb.Point {
	return &pb.Point{
		Time:  timestamppb.New(t),
		Value: v,
	}
}
