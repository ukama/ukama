/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

// Package specsync enforces that the aggregator's KPI spec copy is
// byte-identical to the canonical set in analysis/configs/kpis. The two
// services ship separate images (each COPYies its own configs/), so the
// files exist twice on disk — this test is what makes silent drift a CI
// failure instead of a production inconsistency. Fix a failure with
// `make sync-specs` in aggregator/.
package specsync

import (
	"os"
	"path/filepath"
	"testing"
)

const (
	canonicalDir = "../../../analysis/configs/kpis"
	copyDir      = "../../configs/kpis"
)

func yamlSet(t *testing.T, dir string) map[string][]byte {
	t.Helper()

	files, err := filepath.Glob(filepath.Join(dir, "*.yaml"))
	if err != nil {
		t.Fatalf("globbing %s: %v", dir, err)
	}

	if len(files) == 0 {
		t.Fatalf("no KPI specs found in %s", dir)
	}

	out := map[string][]byte{}

	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			t.Fatalf("reading %s: %v", f, err)
		}

		out[filepath.Base(f)] = b
	}

	return out
}

func TestKpiSpecsMatchCanonical(t *testing.T) {
	canonical := yamlSet(t, canonicalDir)
	copied := yamlSet(t, copyDir)

	for name := range canonical {
		if _, ok := copied[name]; !ok {
			t.Errorf("spec %s exists in analysis/configs/kpis but not here — run `make sync-specs`", name)
		}
	}

	for name, body := range copied {
		want, ok := canonical[name]
		if !ok {
			t.Errorf("spec %s exists here but not in analysis/configs/kpis (stale copy?) — run `make sync-specs`", name)

			continue
		}

		if string(body) != string(want) {
			t.Errorf("spec %s differs from the canonical copy — run `make sync-specs`", name)
		}
	}
}
