/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package schema

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

// Pull strategies.
type Strategy string

const (
	// StrategyWindow pulls window-bounded data ({{.WindowStart}}/{{.WindowEnd}}
	// templated into params); eligible at window close.
	StrategyWindow Strategy = "window"
	// StrategyFullSnapshot pulls current state every window, stored as a
	// change-log (rows only on content change + tombstones); eligible at
	// window end.
	StrategyFullSnapshot Strategy = "full_snapshot"
)

// SourceSpec is one source system: what to pull (never when — cadence is
// pipeline config).
type SourceSpec struct {
	Version int        `yaml:"version"`
	Source  string     `yaml:"source"`
	System  string     `yaml:"system"`   // logical name, resolved via initclient
	BaseURL string     `yaml:"base_url"` // optional override, skips initclient
	Pulls   []PullSpec `yaml:"pulls"`
}

// PullSpec is one dataset. Key is the static dataset key
// (<system>.<resource>.<operation>) — globally unique, referenced by KPI
// specs, stamped on raw_records, tracked in the ledger.
type PullSpec struct {
	Key      string            `yaml:"key"`
	Endpoint string            `yaml:"endpoint"` // may contain {{.bind}} templates
	Strategy Strategy          `yaml:"strategy"`
	Params   map[string]string `yaml:"params"`
	Items    string            `yaml:"items"` // path to result array; "$" = bare array
	// Entity is the mapped field used as the entity key (snapshots). A
	// comma-separated list builds a composite key (field values joined with
	// "|") for sources where no single field identifies the entity — e.g. one
	// Prometheus series per iccid×package×site. A component suffixed with
	// "?" is optional: it contributes an empty component when absent rather
	// than failing the pull.
	Entity    string            `yaml:"entity"`
	ForEach   *ForEachSpec      `yaml:"for_each"`
	Map       map[string]string `yaml:"map"` // field -> $.path into each item
	RateLimit string            `yaml:"rate_limit"`

	// System is inherited from the enclosing SourceSpec unless overridden
	// (a pull may target another system's gateway, e.g. subscriber source
	// iterating registry networks).
	System  string `yaml:"system"`
	BaseURL string `yaml:"base_url"`

	// Gateway selects which of the system's gateways to resolve via
	// initclient: "api" (default) or "node" (device-facing node-gateway).
	Gateway string `yaml:"gateway"`

	// OnError: "fail" (default — the whole dataset window fails and is
	// retried) or "record" — a failed for_each iteration writes a synthetic
	// row carrying the binds plus unreachable:true. Use for health probes
	// where unreachable IS the signal.
	OnError string `yaml:"on_error"`
}

// ForEachSpec fans a pull out over the rows of a parent dataset, binding row
// fields into endpoint/params templates. Bound fields plus the parent's
// lineage are stamped onto every child row.
type ForEachSpec struct {
	Dataset string         `yaml:"dataset"`
	Bind    []string       `yaml:"bind"`
	Filter  *ForEachFilter `yaml:"filter"` // optional parent-row filter
}

// ForEachFilter keeps only parent rows whose field value is in the list
// (case-insensitive), e.g. only tnode/anode nodes for health probes.
type ForEachFilter struct {
	Field string   `yaml:"field"`
	In    []string `yaml:"in"`
}

// KPI kinds: what a KPI *is* determines how it aggregates — callers never
// choose a fold function (the Prometheus counter/gauge pattern).
const (
	// KindFlow is an amount that accrues over time (bytes, sales, cents):
	// SUM over time, SUM across scopes.
	KindFlow = "flow"
	// KindGauge is a level that exists at a moment (customers, MRR, an
	// uptime ratio): latest bucket over time; across scopes per ScopeAgg.
	KindGauge = "gauge"
)

// Scope aggregations for gauges.
const (
	ScopeAggSum = "sum" // additive gauges: customers, MRR, sites online
	ScopeAggAvg = "avg" // ratios: uptime — weighted via components
)

// KpiSpec is one KPI: inputs from the warehouse (never endpoints), one
// dedicated algo, output metadata, and the rollup ops Aggregator may apply.
type KpiSpec struct {
	Kpi    string `yaml:"kpi"`
	Domain string `yaml:"domain"`
	Algo   string `yaml:"algo"` // name@version, from the algo registry
	// Kind (flow|gauge) drives the default aggregation server-side so
	// callers ask for the KPI, not for a fold function.
	Kind string `yaml:"kind"`
	// ScopeAgg (gauges only): how the gauge folds ACROSS scopes — sum for
	// additive gauges, avg (weighted) for ratios. Flows always sum.
	ScopeAgg          string               `yaml:"scope_agg"`
	Scope             []string             `yaml:"scope"`
	Inputs            map[string]InputSpec `yaml:"inputs"`
	Output            OutputSpec           `yaml:"output"`
	PositiveDirection string               `yaml:"positive_direction"` // up|down (console polarity)
	Lookback          string               `yaml:"lookback"`           // optional, e.g. "30d"
	// Params are optional algo-specific tuning knobs, validated by the algo
	// that consumes them.
	Params map[string]string `yaml:"params"`
}

// DefaultReadOp is the aggregation the query planner computes for this KPI
// when the caller does not override it:
// flow → SUM; gauge ratios → AVG (weighted); additive gauges → LAST
// (current level; folded across scopes as a sum of latest values).
// Every op is computable at read time from a rollup row's components
// (Sum/Count/Min/Max/Last) — nothing is materialized per op.
func (k KpiSpec) DefaultReadOp() string {
	if k.Kind == KindGauge {
		if k.ScopeAgg == ScopeAggAvg {
			return "AVG"
		}

		return "LAST"
	}

	return "SUM"
}

// InputSpec references a dataset by its static key.
type InputSpec struct {
	Dataset string   `yaml:"dataset"`
	Fields  []string `yaml:"fields"`
	// Mode: "state" (default; state-as-of-window for snapshot datasets),
	// "window" (only rows belonging to the window), or "state_prev"
	// (state as of the PREVIOUS window — the lag-1 baseline that lets an
	// algo turn a cumulative counter into a per-window increment).
	Mode string `yaml:"mode"`
}

type OutputSpec struct {
	Type   string `yaml:"type"`
	Unit   string `yaml:"unit"`
	Symbol string `yaml:"symbol"`
}

// Read-time aggregations, all computable from a rollup row's components.
var ValidOps = map[string]bool{
	"SUM": true, "AVG": true, "MIN": true, "MAX": true, "COUNT": true, "LAST": true,
}

// LoadSourceSpecs reads and validates all *.yaml files in dir.
func LoadSourceSpecs(dir string) ([]SourceSpec, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.yaml"))
	if err != nil {
		return nil, err
	}

	specs := make([]SourceSpec, 0, len(files))

	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			return nil, fmt.Errorf("reading source spec %s: %w", f, err)
		}

		var s SourceSpec
		if err := yaml.Unmarshal(b, &s); err != nil {
			return nil, fmt.Errorf("parsing source spec %s: %w", f, err)
		}

		for i := range s.Pulls {
			if s.Pulls[i].System == "" {
				s.Pulls[i].System = s.System
			}
			if s.Pulls[i].BaseURL == "" {
				s.Pulls[i].BaseURL = s.BaseURL
			}
			if s.Pulls[i].Items == "" {
				s.Pulls[i].Items = "$"
			}
		}

		specs = append(specs, s)
	}

	if err := ValidateSourceSpecs(specs); err != nil {
		return nil, err
	}

	return specs, nil
}

// ValidateSourceSpecs enforces: global dataset-key uniqueness, valid
// strategies, for_each references to existing keys, snapshot entity keys,
// and DAG acyclicity.
func ValidateSourceSpecs(specs []SourceSpec) error {
	keys := map[string]bool{}

	for _, s := range specs {
		for _, p := range s.Pulls {
			if p.Key == "" {
				return fmt.Errorf("source %s: pull with empty dataset key", s.Source)
			}
			if keys[p.Key] {
				return fmt.Errorf("duplicate dataset key %q — one endpoint, one key, one definition", p.Key)
			}
			keys[p.Key] = true

			if p.Strategy != StrategyWindow && p.Strategy != StrategyFullSnapshot {
				return fmt.Errorf("dataset %s: unknown strategy %q", p.Key, p.Strategy)
			}
			if p.Strategy == StrategyFullSnapshot && p.Entity == "" {
				return fmt.Errorf("dataset %s: full_snapshot requires an entity field", p.Key)
			}
			for _, ef := range p.EntityFields() {
				if _, mapped := p.Map[ef.Name]; mapped {
					continue
				}

				bound := false
				if p.ForEach != nil {
					for _, b := range p.ForEach.Bind {
						if b == ef.Name {
							bound = true

							break
						}
					}
				}

				if !bound {
					return fmt.Errorf("dataset %s: entity field %q is neither mapped nor bound", p.Key, ef.Name)
				}
			}
			if p.Gateway != "" && p.Gateway != "api" && p.Gateway != "node" {
				return fmt.Errorf("dataset %s: gateway must be api or node, got %q", p.Key, p.Gateway)
			}
			if p.OnError != "" && p.OnError != "fail" && p.OnError != "record" {
				return fmt.Errorf("dataset %s: on_error must be fail or record, got %q", p.Key, p.OnError)
			}
			if p.OnError == "record" && p.ForEach == nil {
				return fmt.Errorf("dataset %s: on_error record requires for_each (the binds identify the failed entity)", p.Key)
			}
			if len(p.Map) == 0 {
				return fmt.Errorf("dataset %s: empty field map", p.Key)
			}
		}
	}

	for _, s := range specs {
		for _, p := range s.Pulls {
			if p.ForEach != nil && !keys[p.ForEach.Dataset] {
				return fmt.Errorf("dataset %s: for_each references unknown dataset %q", p.Key, p.ForEach.Dataset)
			}
		}
	}

	if _, err := OrderPulls(specs); err != nil {
		return err
	}

	return nil
}

// OrderPulls topologically sorts all pulls into execution stages: stage 0 has
// no dependencies; stage k depends only on earlier stages. Errors on cycles.
func OrderPulls(specs []SourceSpec) ([][]PullSpec, error) {
	pulls := map[string]PullSpec{}
	for _, s := range specs {
		for _, p := range s.Pulls {
			pulls[p.Key] = p
		}
	}

	depth := map[string]int{}

	var resolve func(key string, seen map[string]bool) (int, error)
	resolve = func(key string, seen map[string]bool) (int, error) {
		if d, ok := depth[key]; ok {
			return d, nil
		}
		if seen[key] {
			return 0, fmt.Errorf("for_each cycle involving dataset %q", key)
		}
		seen[key] = true

		p := pulls[key]
		d := 0
		if p.ForEach != nil {
			pd, err := resolve(p.ForEach.Dataset, seen)
			if err != nil {
				return 0, err
			}
			d = pd + 1
		}

		depth[key] = d

		return d, nil
	}

	maxDepth := 0
	for key := range pulls {
		d, err := resolve(key, map[string]bool{})
		if err != nil {
			return nil, err
		}
		if d > maxDepth {
			maxDepth = d
		}
	}

	stages := make([][]PullSpec, maxDepth+1)
	// Deterministic order within a stage.
	keys := make([]string, 0, len(pulls))
	for k := range pulls {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		stages[depth[k]] = append(stages[depth[k]], pulls[k])
	}

	return stages, nil
}

// LoadKpiSpecs reads and validates all *.yaml files in dir.
func LoadKpiSpecs(dir string) ([]KpiSpec, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.yaml"))
	if err != nil {
		return nil, err
	}

	specs := make([]KpiSpec, 0, len(files))
	seen := map[string]bool{}

	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			return nil, fmt.Errorf("reading kpi spec %s: %w", f, err)
		}

		var s KpiSpec
		if err := yaml.Unmarshal(b, &s); err != nil {
			return nil, fmt.Errorf("parsing kpi spec %s: %w", f, err)
		}

		if s.Kpi == "" || s.Algo == "" {
			return nil, fmt.Errorf("kpi spec %s: kpi and algo are required", f)
		}
		if s.Kind != KindFlow && s.Kind != KindGauge {
			return nil, fmt.Errorf("kpi %s: kind must be flow or gauge, got %q", s.Kpi, s.Kind)
		}
		if s.Kind == KindGauge && s.ScopeAgg != ScopeAggSum && s.ScopeAgg != ScopeAggAvg {
			return nil, fmt.Errorf("kpi %s: gauge requires scope_agg sum or avg, got %q", s.Kpi, s.ScopeAgg)
		}
		if s.Kind == KindFlow && s.ScopeAgg != "" && s.ScopeAgg != ScopeAggSum {
			return nil, fmt.Errorf("kpi %s: flow scope_agg is always sum (omit it), got %q", s.Kpi, s.ScopeAgg)
		}
		if seen[s.Kpi] {
			return nil, fmt.Errorf("duplicate kpi key %q", s.Kpi)
		}
		seen[s.Kpi] = true

		if len(s.Inputs) == 0 {
			return nil, fmt.Errorf("kpi %s: at least one input dataset required", s.Kpi)
		}
		for name, in := range s.Inputs {
			if in.Dataset == "" {
				return nil, fmt.Errorf("kpi %s: input %s has no dataset key", s.Kpi, name)
			}
		}

		specs = append(specs, s)
	}

	return specs, nil
}

// EntityField is one component of a (possibly composite) entity key.
type EntityField struct {
	Name string
	// Optional components (declared with a trailing "?") may be absent from
	// a row; they contribute an empty component instead of failing the pull.
	Optional bool
}

// EntityFields returns the entity key's component fields (one for a simple
// key, several for a composite "a,b,c" key).
func (p PullSpec) EntityFields() []EntityField {
	if p.Entity == "" {
		return nil
	}

	parts := strings.Split(p.Entity, ",")
	out := make([]EntityField, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		field := EntityField{Name: part}
		if strings.HasSuffix(part, "?") {
			field.Name = strings.TrimSuffix(part, "?")
			field.Optional = true
		}

		out = append(out, field)
	}

	return out
}

// InputDatasets returns the set of dataset keys a KPI depends on.
func (k KpiSpec) InputDatasets() []string {
	out := make([]string, 0, len(k.Inputs))
	for _, in := range k.Inputs {
		out = append(out, in.Dataset)
	}
	sort.Strings(out)

	return out
}
