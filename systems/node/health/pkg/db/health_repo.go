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
	"gorm.io/gorm/clause"
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
			q = q.Where(map[string]interface{}{"nodeId": nodeID})
		}
		if reportID != "" {
			q = q.Where(map[string]interface{}{"reportId": reportID})
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

	q := r.Db.GetGormDb().Model(&HealthReport{}).Order(clause.OrderByColumn{
		Column: clause.Column{Name: "reportedAt"},
		Desc:   true,
	})
	if reportID != "" {
		q = q.Where("id = ?", reportID)
	}
	if nodeID != "" {
		q = q.Where(map[string]interface{}{"nodeId": nodeID})
	}
	if reportedAt != nil {
		q = q.Where(map[string]interface{}{"reportedAt": *reportedAt})
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
	err := tx.Where(map[string]interface{}{"nodeId": report.NodeID}).First(&existing).Error
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
	return tx.Model(&Node{}).Where(map[string]interface{}{"nodeId": report.NodeID}).Updates(map[string]interface{}{
		"lastSeenAt":     receivedAt,
		"nodeType":       report.NodeType,
		"lastReportedAt": report.ReportedAt,
	}).Error
}

func syncNodeLatestHealth(tx *gorm.DB, report *HealthReport, receivedAt time.Time) error {
	var current NodeLatestHealth
	err := tx.Where(map[string]interface{}{"nodeId": report.NodeID}).First(&current).Error
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
	return tx.Model(&NodeLatestHealth{}).Where(map[string]interface{}{"nodeId": report.NodeID}).Updates(map[string]interface{}{
		"nodeType":      next.NodeType,
		"reportId":      next.ReportID,
		"schemaVersion": next.SchemaVersion,
		"reportedAt":    next.ReportedAt,
		"receivedAt":    next.ReceivedAt,
		"parseStatus":   next.ParseStatus,
		"parseError":    next.ParseError,
		"payload":       next.Payload,
		"updatedAt":     next.UpdatedAt,
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
