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
	Key       string            `yaml:"key"`
	Endpoint  string            `yaml:"endpoint"` // may contain {{.bind}} templates
	Strategy  Strategy          `yaml:"strategy"`
	Params    map[string]string `yaml:"params"`
	Items     string            `yaml:"items"`  // path to result array; "$" = bare array
	Entity    string            `yaml:"entity"` // mapped field used as entity key (snapshots)
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

// KpiSpec is one KPI: inputs from the warehouse (never endpoints), one
// dedicated algo, output metadata, and the rollup ops Aggregator may apply.
type KpiSpec struct {
	Kpi               string               `yaml:"kpi"`
	Domain            string               `yaml:"domain"`
	Algo              string               `yaml:"algo"` // name@version, from the algo registry
	Scope             []string             `yaml:"scope"`
	Inputs            map[string]InputSpec `yaml:"inputs"`
	Output            OutputSpec           `yaml:"output"`
	RollupOps         []string             `yaml:"rollup_ops"`
	PositiveDirection string               `yaml:"positive_direction"` // up|down (console polarity)
	Lookback          string               `yaml:"lookback"`           // optional, e.g. "30d"
}

// InputSpec references a dataset by its static key.
type InputSpec struct {
	Dataset string   `yaml:"dataset"`
	Fields  []string `yaml:"fields"`
	// Mode: "state" (default; state-as-of-window for snapshot datasets) or
	// "window" (only rows belonging to the window).
	Mode string `yaml:"mode"`
}

type OutputSpec struct {
	Type   string `yaml:"type"`
	Unit   string `yaml:"unit"`
	Symbol string `yaml:"symbol"`
}

// Rollup operations.
var ValidOps = map[string]bool{
	"SUM": true, "AVG": true, "MIN": true, "MAX": true, "COUNT": true, "LAST": true, "DELTA": true,
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
		if len(s.RollupOps) == 0 {
			return nil, fmt.Errorf("kpi %s: rollup_ops required", s.Kpi)
		}
		for _, op := range s.RollupOps {
			if !ValidOps[strings.ToUpper(op)] {
				return nil, fmt.Errorf("kpi %s: invalid rollup op %q", s.Kpi, op)
			}
		}

		specs = append(specs, s)
	}

	return specs, nil
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
