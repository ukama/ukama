/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"errors"
	"time"

	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/ukama"
	"gorm.io/gorm"
)

type HealthRepo interface {
	List(reportID, nodeID string, reportedAt *time.Time, timeframe ukama.FilterTimeframesType) ([]*HealthReport, error)
	StoreHealthReport(report *HealthReport, receivedAt time.Time) error
}

type healthRepo struct {
	Db sql.Db
}

func NewHealthRepo(db sql.Db) HealthRepo {
	return &healthRepo{
		Db: db,
	}
}

func latestToReport(l *NodeLatestHealth) *HealthReport {
	return &HealthReport{
		ID:            l.ReportID,
		NodeID:        l.NodeID,
		NodeType:      l.NodeType,
		SchemaVersion: l.SchemaVersion,
		ReportedAt:    l.ReportedAt,
		ReceivedAt:    l.ReceivedAt,
		ParseStatus:   l.ParseStatus,
		ParseError:    l.ParseError,
		Payload:       l.Payload,
	}
}

func (r *healthRepo) List(reportID, nodeID string, reportedAt *time.Time, timeframe ukama.FilterTimeframesType) ([]*HealthReport, error) {
	if timeframe == ukama.FilterTimeframesTypeLatest {
		q := r.Db.GetGormDb().Model(&NodeLatestHealth{})
		if nodeID != "" {
			q = q.Where("node_id = ?", nodeID)
		}
		if reportID != "" {
			q = q.Where("report_id = ?", reportID)
		}
		var rows []NodeLatestHealth
		if err := q.Find(&rows).Error; err != nil {
			return nil, err
		}
		out := make([]*HealthReport, len(rows))
		for i := range rows {
			out[i] = latestToReport(&rows[i])
		}
		return out, nil
	}

	q := r.Db.GetGormDb().Model(&HealthReport{}).Order("reported_at DESC")
	if reportID != "" {
		q = q.Where("id = ?", reportID)
	}
	if nodeID != "" {
		q = q.Where("node_id = ?", nodeID)
	}
	if reportedAt != nil {
		q = q.Where("reported_at = ?", *reportedAt)
	}
	var reports []*HealthReport
	err := q.Find(&reports).Error
	return reports, err
}

func shouldReplaceNodeLatest(report *HealthReport, current *NodeLatestHealth) bool {
	if report.ReportedAt.After(current.ReportedAt) {
		return true
	}
	return report.ReportedAt.Equal(current.ReportedAt) && report.ReceivedAt.After(current.ReceivedAt)
}

func upsertHealthNode(tx *gorm.DB, report *HealthReport, receivedAt time.Time) error {
	var existing Node
	err := tx.Where("node_id = ?", report.NodeID).First(&existing).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		node := &Node{
			NodeID:         report.NodeID,
			NodeType:       report.NodeType,
			FirstSeenAt:    receivedAt,
			LastSeenAt:     receivedAt,
			LastReportedAt: &report.ReportedAt,
		}
		return tx.Create(node).Error
	}
	return tx.Model(&Node{}).Where("node_id = ?", report.NodeID).Updates(map[string]interface{}{
		"last_seen_at":      receivedAt,
		"node_type":         report.NodeType,
		"last_reported_at": report.ReportedAt,
	}).Error
}

func syncNodeLatestHealth(tx *gorm.DB, report *HealthReport, receivedAt time.Time) error {
	var current NodeLatestHealth
	err := tx.Where("node_id = ?", report.NodeID).First(&current).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	next := NodeLatestHealth{
		NodeID:        report.NodeID,
		NodeType:      report.NodeType,
		ReportID:      report.ID,
		SchemaVersion: report.SchemaVersion,
		ReportedAt:    report.ReportedAt,
		ReceivedAt:    report.ReceivedAt,
		ParseStatus:   report.ParseStatus,
		ParseError:    report.ParseError,
		Payload:       report.Payload,
		UpdatedAt:     receivedAt,
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return tx.Create(&next).Error
	}
	if !shouldReplaceNodeLatest(report, &current) {
		return nil
	}
	return tx.Model(&NodeLatestHealth{}).Where("node_id = ?", report.NodeID).Updates(map[string]interface{}{
		"node_type":      next.NodeType,
		"report_id":      next.ReportID,
		"schema_version": next.SchemaVersion,
		"reported_at":    next.ReportedAt,
		"received_at":    next.ReceivedAt,
		"parse_status":   next.ParseStatus,
		"parse_error":    next.ParseError,
		"payload":        next.Payload,
		"updated_at":     next.UpdatedAt,
	}).Error
}

func (r *healthRepo) StoreHealthReport(report *HealthReport, receivedAt time.Time) error {
	if report == nil {
		return errors.New("nil health report")
	}
	report.ReceivedAt = receivedAt

	return r.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(report).Error; err != nil {
			return err
		}
		if err := upsertHealthNode(tx, report, receivedAt); err != nil {
			return err
		}
		return syncNodeLatestHealth(tx, report, receivedAt)
	})
}
