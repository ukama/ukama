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
	"strings"

	"gopkg.in/yaml.v2"
)

// ReportSpec declares a resource performance report: entity rows from a raw
// dataset's state, KPI columns from existing KPIs, threshold-based status
// labels. Config wires, code computes.
type ReportSpec struct {
	Report   string          `yaml:"report"`
	Title    string          `yaml:"title"`
	Resource ReportResource  `yaml:"resource"`
	RowScope string          `yaml:"row_scope"` // KPI scope dimension keyed by entity
	Columns  []ReportColumn  `yaml:"columns"`
	Status   []StatusRule    `yaml:"status"`
	Sort     ReportSort      `yaml:"sort"`
}

type ReportResource struct {
	Dataset      string            `yaml:"dataset"`       // entity rows = state-as-of latest
	EntityKey    string            `yaml:"entity_key"`    // field holding the entity id
	NetworkMatch string            `yaml:"network_match"` // entity field matched against network filter; "" on the entity = org-level (matches all)
	Attributes   []ReportAttribute `yaml:"attributes"`
}

type ReportAttribute struct {
	Name   string `yaml:"name"`
	Field  string `yaml:"field"`
	Format string `yaml:"format"` // money|days|bytes|"" (display hint)
}

type ReportColumn struct {
	Name   string `yaml:"name"`
	Kpi    string `yaml:"kpi"`
	Op     string `yaml:"op"`
	Format string `yaml:"format"`
}

// StatusRule: first matching rule wins. When is a minimal comparison
// "<column-or-attribute> <op> <literal>" with op in == != < <= > >=.
type StatusRule struct {
	When    string `yaml:"when"`
	Default bool   `yaml:"default"`
	Label   string `yaml:"label"`
}

type ReportSort struct {
	Column string `yaml:"column"`
	Desc   bool   `yaml:"desc"`
}

// LoadReportSpecs reads and validates all *.yaml files in dir against the
// KPI registry: every column must reference an existing KPI, an allowed
// rollup op, and a KPI scoped by the report's row_scope.
func LoadReportSpecs(dir string, kpis []KpiSpec) ([]ReportSpec, error) {
	byKey := map[string]KpiSpec{}
	for _, k := range kpis {
		byKey[k.Kpi] = k
	}

	files, err := filepath.Glob(filepath.Join(dir, "*.yaml"))
	if err != nil {
		return nil, err
	}

	specs := make([]ReportSpec, 0, len(files))
	seen := map[string]bool{}

	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			return nil, fmt.Errorf("reading report spec %s: %w", f, err)
		}

		var s ReportSpec
		if err := yaml.Unmarshal(b, &s); err != nil {
			return nil, fmt.Errorf("parsing report spec %s: %w", f, err)
		}

		if s.Report == "" || s.Resource.Dataset == "" || s.Resource.EntityKey == "" || s.RowScope == "" {
			return nil, fmt.Errorf("report spec %s: report, resource.dataset, resource.entity_key and row_scope are required", f)
		}

		if seen[s.Report] {
			return nil, fmt.Errorf("duplicate report %q", s.Report)
		}
		seen[s.Report] = true

		if len(s.Columns) == 0 {
			return nil, fmt.Errorf("report %s: at least one column required", s.Report)
		}

		for i, col := range s.Columns {
			kpi, ok := byKey[col.Kpi]
			if !ok {
				return nil, fmt.Errorf("report %s column %s: unknown kpi %q", s.Report, col.Name, col.Kpi)
			}

			op := strings.ToUpper(col.Op)
			if op == "" {
				op = kpi.DefaultReadOp() // kind-derived, like the query API
			}
			s.Columns[i].Op = op

			if !ValidOps[op] {
				return nil, fmt.Errorf("report %s column %s: unknown op %s (valid: SUM, AVG, MIN, MAX, COUNT, LAST)",
					s.Report, col.Name, op)
			}

			scoped := false
			for _, sc := range kpi.Scope {
				if sc == s.RowScope {
					scoped = true

					break
				}
			}

			if !scoped {
				return nil, fmt.Errorf("report %s column %s: kpi %s is not scoped by %s", s.Report, col.Name, col.Kpi, s.RowScope)
			}
		}

		for _, rule := range s.Status {
			if !rule.Default && rule.When == "" {
				return nil, fmt.Errorf("report %s: status rule needs 'when' or 'default: true'", s.Report)
			}

			if rule.Label == "" {
				return nil, fmt.Errorf("report %s: status rule missing label", s.Report)
			}
		}

		specs = append(specs, s)
	}

	return specs, nil
}
