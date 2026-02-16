/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package algos

import (
	"testing"
)

func TestEvaluateDomain_NoMatchingRules_ReturnsHealthy(t *testing.T) {
	rules := []Rule{
		{ID: "health.cpu.pressure", Domain: "health", Severity: "degraded", Conditions: []Condition{
			{Metric: "cpu", State: "warning"},
		}, EvidenceMetrics: []string{"cpu"}},
	}
	evals := MetricEvaluationsMap{
		"cpu":    {MetricID: "cpu", State: "healthy", Trend: "stable", Confidence: 0.9},
		"memory": {MetricID: "memory", State: "healthy", Trend: "stable", Confidence: 0.9},
	}

	snap := EvaluateDomain("health", evals, rules, nil, 1000)
	if snap.RuleID != "health.healthy" {
		t.Errorf("expected health.healthy, got %s", snap.RuleID)
	}
	if snap.Severity != "healthy" {
		t.Errorf("expected severity healthy, got %s", snap.Severity)
	}
}

func TestEvaluateDomain_MatchingRule_ReturnsRule(t *testing.T) {
	rules := []Rule{
		{ID: "health.cpu.pressure", Domain: "health", Severity: "degraded", Priority: 60, Conditions: []Condition{
			{Metric: "cpu", State: "warning"},
			{Metric: "cpu", Trend: "increasing"},
			{Metric: "cpu", ConfidenceGte: 0.7},
		}, EvidenceMetrics: []string{"cpu"}},
	}
	evals := MetricEvaluationsMap{
		"cpu":    {MetricID: "cpu", State: "warning", Trend: "increasing", Confidence: 0.8},
		"memory": {MetricID: "memory", State: "healthy", Trend: "stable", Confidence: 0.9},
	}

	snap := EvaluateDomain("health", evals, rules, nil, 1000)
	if snap.RuleID != "health.cpu.pressure" {
		t.Errorf("expected health.cpu.pressure, got %s", snap.RuleID)
	}
	if snap.Severity != "degraded" {
		t.Errorf("expected severity degraded, got %s", snap.Severity)
	}
}

func TestEvaluateDomain_EscalateImmediately(t *testing.T) {
	rules := []Rule{
		{ID: "health.cpu.pressure", Domain: "health", Severity: "degraded", Conditions: []Condition{
			{Metric: "cpu", State: "warning"},
		}, EvidenceMetrics: []string{"cpu"}},
		{ID: "health.cpu.critical", Domain: "health", Severity: "critical", Conditions: []Condition{
			{Metric: "cpu", State: "critical"},
		}, EvidenceMetrics: []string{"cpu"}},
	}
	evals := MetricEvaluationsMap{
		"cpu": {MetricID: "cpu", State: "critical", Trend: "increasing", Confidence: 0.8},
	}
	previous := &DomainSnapshot{RuleID: "health.cpu.pressure", Severity: "degraded", Domain: "health"}

	snap := EvaluateDomain("health", evals, rules, previous, 1000)
	if snap.RuleID != "health.cpu.critical" {
		t.Errorf("expected immediate escalation to health.cpu.critical, got %s", snap.RuleID)
	}
}

func TestLoadRulesFromJSON(t *testing.T) {
	rules, err := LoadRulesFromJSON("../../metric-rules.json")
	if err != nil {
		t.Skipf("metric-rules.json not found: %v", err)
	}
	if len(rules) == 0 {
		t.Error("expected at least one rule")
	}
	if rules[0].Domain != "health" {
		t.Errorf("expected domain health, got %s", rules[0].Domain)
	}
}

func TestLoadRulesFromJSONAnyOfParsing(t *testing.T) {
	rules, err := LoadRulesFromJSON("../../metric-rules.json")
	if err != nil {
		t.Skipf("metric-rules.json not found: %v", err)
	}
	var resRule *Rule
	for i := range rules {
		if rules[i].ID == "health.resource.constraint" {
			resRule = &rules[i]
			break
		}
	}
	if resRule == nil {
		t.Fatal("health.resource.constraint rule not found")
	}
	if len(resRule.Conditions) < 2 {
		t.Fatalf("expected 2 conditions, got %d", len(resRule.Conditions))
	}
	// First condition: any_of with cpu/warning and memory/warning, min_match 2
	if len(resRule.Conditions[0].AnyOf) != 2 {
		t.Errorf("expected 2 any_of items in condition 0, got %d", len(resRule.Conditions[0].AnyOf))
	}
	if resRule.Conditions[0].AnyOf[0].Metric != "cpu" || resRule.Conditions[0].AnyOf[0].State != "warning" {
		t.Errorf("any_of[0] expected cpu/warning, got %q/%q", resRule.Conditions[0].AnyOf[0].Metric, resRule.Conditions[0].AnyOf[0].State)
	}
	if resRule.Conditions[0].AnyOf[1].Metric != "memory" || resRule.Conditions[0].AnyOf[1].State != "warning" {
		t.Errorf("any_of[1] expected memory/warning, got %q/%q", resRule.Conditions[0].AnyOf[1].Metric, resRule.Conditions[0].AnyOf[1].State)
	}
	// Second condition: any_of with trend increasing
	if len(resRule.Conditions[1].AnyOf) != 2 {
		t.Errorf("expected 2 any_of items in condition 1, got %d", len(resRule.Conditions[1].AnyOf))
	}
	if resRule.Conditions[1].AnyOf[0].Trend != "increasing" || resRule.Conditions[1].AnyOf[1].Trend != "increasing" {
		t.Errorf("any_of expected trend increasing, got %q and %q", resRule.Conditions[1].AnyOf[0].Trend, resRule.Conditions[1].AnyOf[1].Trend)
	}
}
