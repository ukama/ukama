/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package schema

import (
	"time"
)

// Zone ownership (enforced by service wiring, later by DB roles):
//   ingest     writes RawRecord, WindowLedger(dataset), IngestError
//   analysis   writes KpiWindow, WindowLedger(kpi), AnalysisError; reads raw zone
//   aggregator writes KpiRollup; reads kpi zone

// Ledger kinds.
const (
	LedgerKindDataset = "dataset"
	LedgerKindKpi     = "kpi"
)

// Ledger statuses (state machine: in_flight -> pulled/computed | failed; dirty
// re-enters the pipeline at the owning stage).
const (
	StatusInFlight = "in_flight"
	StatusPulled   = "pulled"
	StatusComputed = "computed"
	StatusFailed   = "failed"
	StatusDirty    = "dirty"
)

// RawRecord is the bronze zone: one row per source record per window.
// Payload is the untouched source JSON (replayable); Fields is the
// spec-mapped normalized projection.
//
// Windowed pulls: one row per source record, deduped on DedupKey.
// Snapshot pulls (change-log): one row per entity *only when its content
// hash changes*, plus tombstones (Deleted=true) when an entity disappears.
// "State as of window N" = latest row per EntityKey with WindowID <= N,
// excluding tombstones.
type RawRecord struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement"`
	OrgID       string    `gorm:"size:64;index:idx_raw_dsw,priority:1;uniqueIndex:uq_raw_dedup,priority:1"`
	DatasetKey  string    `gorm:"size:128;index:idx_raw_dsw,priority:2;uniqueIndex:uq_raw_dedup,priority:2"`
	WindowID    int64     `gorm:"index:idx_raw_dsw,priority:3;uniqueIndex:uq_raw_dedup,priority:3"`
	EntityKey   string    `gorm:"size:128;index"`
	ContentHash string    `gorm:"size:64"`
	DedupKey    string    `gorm:"size:160;uniqueIndex:uq_raw_dedup,priority:4"`
	Deleted     bool      `gorm:"default:false"`
	EventTime   time.Time `gorm:"index"`
	Payload     string    `gorm:"type:jsonb"`
	Fields      string    `gorm:"type:jsonb"`
	IngestedAt  time.Time
}

// WindowLedger is the pipeline's source of truth for what happened to every
// (kind, key, org, window). msgbus events are only the fast path; sweepers
// recover anything whose event was lost.
type WindowLedger struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement"`
	OrgID     string `gorm:"size:64;uniqueIndex:uq_ledger,priority:1"`
	Kind      string `gorm:"size:16;uniqueIndex:uq_ledger,priority:2"`
	Key       string `gorm:"size:128;uniqueIndex:uq_ledger,priority:3"`
	WindowID  int64  `gorm:"uniqueIndex:uq_ledger,priority:4;index"`
	Status    string `gorm:"size:16;index"`
	Attempt   int
	Detail    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// KpiWindow is the KPI zone: one KPI value per (kpi, org, scope, window).
// Components (Sum/Count/Min/Max) are mandatory — they are what make
// higher-span rollups exact (weighted, never avg-of-avg).
type KpiWindow struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement"`
	KpiKey      string `gorm:"size:64;uniqueIndex:uq_kpiwin,priority:1;index"`
	OrgID       string `gorm:"size:64;uniqueIndex:uq_kpiwin,priority:2"`
	Scope       string `gorm:"size:512;uniqueIndex:uq_kpiwin,priority:3"` // canonical JSON, sorted keys
	WindowID    int64  `gorm:"uniqueIndex:uq_kpiwin,priority:4;index"`
	Value       float64
	Sum         float64
	Count       float64
	Min         float64
	Max         float64
	ValueType   string `gorm:"size:16"`
	Unit        string `gorm:"size:16"`
	Symbol      string `gorm:"size:16"`
	AlgoVersion string `gorm:"size:64"`
	ComputedAt  time.Time
}

// RollupRowOp marks the single components row the engine materializes per
// (kpi, org, scope, span, span_start). Historically one row existed PER OP
// (SUM/AVG/... with a precomputed Value); every aggregation is now computed
// at read time from the components, so exactly one row is written, tagged
// with this marker. The column stays in the unique index so old per-op rows
// coexist harmlessly (reads filter on the marker; the boot backfill
// re-materializes history).
const RollupRowOp = "VAL"

// KpiRollup is the rollup zone: ONE components row per (kpi, org, scope,
// span, span_start).
//
// Sum/Count/Min/Max/Last are the span's aggregated components (folded from
// the span's kpi_windows): any read-time aggregation — including exact
// weighted AVG and cross-scope group_by folds — derives from them. Value
// caches the KPI's kind-default aggregation for debugging/BI convenience.
// Trend is computed at read time (same question over the previous period),
// not stored.
type KpiRollup struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement"`
	KpiKey     string    `gorm:"size:64;uniqueIndex:uq_rollup,priority:1;index"`
	OrgID      string    `gorm:"size:64;uniqueIndex:uq_rollup,priority:2"`
	Scope      string    `gorm:"size:512;uniqueIndex:uq_rollup,priority:3"`
	Span       string    `gorm:"size:16;uniqueIndex:uq_rollup,priority:4"` // daily|weekly|monthly
	SpanStart  time.Time `gorm:"uniqueIndex:uq_rollup,priority:5;index"`
	SpanEnd    time.Time
	Op         string `gorm:"size:16;uniqueIndex:uq_rollup,priority:6"` // always RollupRowOp
	Value      float64
	Sum        float64
	Count      float64
	Min        float64
	Max        float64
	Last       float64 // latest window's value in the span (gauge level)
	ValueType  string  `gorm:"size:16"`
	Unit       string  `gorm:"size:16"`
	Symbol     string  `gorm:"size:16"`
	IsPartial  bool
	ComputedAt time.Time
}

// IngestError records a failed pull (whole dataset window or a single
// for_each iteration). Never silently skipped.
type IngestError struct {
	ID         uint64 `gorm:"primaryKey;autoIncrement"`
	OrgID      string `gorm:"size:64;index"`
	DatasetKey string `gorm:"size:128;index"`
	WindowID   int64  `gorm:"index"`
	Iteration  string // bound params of the failed for_each iteration, if any
	Error      string
	CreatedAt  time.Time
}

// AnalysisError records a failed KPI computation.
type AnalysisError struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement"`
	OrgID     string `gorm:"size:64;index"`
	KpiKey    string `gorm:"size:64;index"`
	WindowID  int64  `gorm:"index"`
	Error     string
	CreatedAt time.Time
}
