/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package algos

import (
	"fmt"

	"github.com/ukama/ukama/systems/analytics/schema"
)

// Datasets maps a KPI spec input name to its rows (each row = the parsed
// spec-mapped Fields of a raw record).
type Datasets map[string][]map[string]interface{}

// Result is one KPI value for one scope combination. Components
// (Sum/Count/Min/Max) are mandatory — Aggregator computes exact weighted
// rollups from them.
type Result struct {
	Scope map[string]string
	Value float64
	Sum   float64
	Count float64
	Min   float64
	Max   float64
}

// Algo is a pure, deterministic function: same window + same inputs =>
// same results, on any host, after any replay.
type Algo func(win schema.Window, in Datasets, spec schema.KpiSpec) ([]Result, error)

// Registry maps "name@version" (the KPI spec's algo field) to its
// implementation. Versioned: algo changes ship as new versions.
type Registry struct {
	algos map[string]Algo
}

func NewRegistry() *Registry {
	return &Registry{algos: map[string]Algo{}}
}

func (r *Registry) Register(nameVersion string, algo Algo) {
	r.algos[nameVersion] = algo
}

func (r *Registry) Get(nameVersion string) (Algo, error) {
	algo, ok := r.algos[nameVersion]
	if !ok {
		return nil, fmt.Errorf("algo %q not registered", nameVersion)
	}

	return algo, nil
}

// Default returns the registry with all production algos registered. A KPI
// spec referencing an unregistered algo fails at startup, not at runtime.
func Default() *Registry {
	r := NewRegistry()

	// v3: site online = ALL its nodes are online (v2: only its cnode).
	r.Register("sites_online@v3", SitesOnline)
	// v2: up = service (tnode cellular) and radio available; a site with no
	// tnode/anode is down rather than silently absent. No planned/unplanned
	// split — every window counts.
	r.Register("site_uptime@v2", SiteUptime)
	r.Register("network_uptime@v2", NetworkUptime)
	// v3: reads the shared all-sims dataset (subscriber.sim.list) and
	// filters status=active in-algo — one sims pull feeds all sim KPIs.
	r.Register("active_customers@v3", ActiveCustomers)
	// v1: per-series increments from the metrics system's cumulative
	// data_usage counter.
	r.Register("data_usage@v1", DataUsage)

	// Package KPIs (see docs/packages-kpi-plan.md).
	r.Register("package_sales@v1", PackageSales)
	r.Register("package_revenue@v1", PackageRevenue)
	r.Register("data_sold@v1", DataSold)
	r.Register("mrr@v1", Mrr)
	r.Register("arpu@v1", Arpu)
	r.Register("customers_on_plan@v1", CustomersOnPlan)
	r.Register("active_plans@v1", ActivePlans)

	// Revenue KPIs (settled payments; org-wide scope).
	r.Register("revenue@v1", Revenue)
	r.Register("paid_customers@v1", PaidCustomers)

	return r
}

// CountResult builds a Result for simple gauge/count KPIs where the window
// value is a single observation (Sum=value, Count=1, Min=Max=value).
func CountResult(scope map[string]string, value float64) Result {
	return Result{
		Scope: scope,
		Value: value,
		Sum:   value,
		Count: 1,
		Min:   value,
		Max:   value,
	}
}
