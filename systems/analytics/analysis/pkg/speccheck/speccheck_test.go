/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

// Package speccheck cross-validates the KPI specs against the ingest source
// specs: every KPI input must reference a dataset key some source actually
// pulls. Without this, a renamed or removed dataset key makes the KPI's
// inputsReady check wait forever and the KPI silently never computes
// (the PACKAGE_DATA_USED failure mode) — here it is a test failure instead.
package speccheck

import (
	"testing"

	"github.com/ukama/ukama/systems/analytics/schema"
)

const (
	kpiDir    = "../../configs/kpis"
	sourceDir = "../../../ingest/configs/sources"
)

func TestKpiInputsExistAsSourceDatasets(t *testing.T) {
	kpis, err := schema.LoadKpiSpecs(kpiDir)
	if err != nil {
		t.Fatalf("loading KPI specs: %v", err)
	}

	sources, err := schema.LoadSourceSpecs(sourceDir)
	if err != nil {
		t.Fatalf("loading source specs: %v", err)
	}

	datasets := map[string]bool{}
	for _, s := range sources {
		for _, p := range s.Pulls {
			datasets[p.Key] = true
		}
	}

	for _, kpi := range kpis {
		for name, in := range kpi.Inputs {
			if !datasets[in.Dataset] {
				t.Errorf("kpi %s input %q references dataset %q, which no source spec pulls — the KPI would silently never compute",
					kpi.Kpi, name, in.Dataset)
			}
		}
	}
}

// TestKpiSpecsLoad is the plain load/validation gate (kind, scope_agg,
// default-read-op materialization) so a bad spec fails CI, not service boot.
func TestKpiSpecsLoad(t *testing.T) {
	kpis, err := schema.LoadKpiSpecs(kpiDir)
	if err != nil {
		t.Fatalf("KPI specs invalid: %v", err)
	}

	if len(kpis) == 0 {
		t.Fatal("no KPI specs found")
	}
}
