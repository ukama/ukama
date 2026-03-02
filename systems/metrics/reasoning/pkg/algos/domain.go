/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package algos

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
)

// RulesFile is the top-level structure of metric-rules.json
type RulesFile struct {
	Version  string `json:"version" yaml:"version"`
	Rules    []Rule `json:"rules" yaml:"rules"`
}

// LoadRulesFromJSON loads rules from a JSON file (metric-rules.json format)
func LoadRulesFromJSON(path string) ([]Rule, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var file RulesFile
	if err := json.Unmarshal(bytes, &file); err != nil {
		return nil, err
	}
	if len(file.Rules) == 0 {
		return nil, errors.New("rules file has no rules")
	}
	return file.Rules, nil
}

const (
	MinStableSec = 60 // De-escalation: require rule to match for this long before switching
)

// Rule represents a domain rule from metric-rules.json
type Rule struct {
	ID                string      `json:"id" yaml:"id"`
	Domain            string      `json:"domain" yaml:"domain"`
	Severity          string      `json:"severity" yaml:"severity"`
	Priority          int         `json:"priority" yaml:"priority"`
	Headline          string      `json:"headline" yaml:"headline"`
	RootCause         string      `json:"root_cause" yaml:"root_cause"`
	ServiceImpact    string      `json:"service_impact" yaml:"service_impact"`
	Conditions        []Condition `json:"conditions" yaml:"conditions"`
	ConditionsMinMatch int        `json:"conditions_min_match" yaml:"conditions_min_match"`
	EvidenceMetrics   []string    `json:"evidence_metrics" yaml:"evidence_metrics"`
}

// Condition is a single rule condition (simple or any_of)
type Condition struct {
	Metric       string   `json:"metric" yaml:"metric"`
	State        string   `json:"state" yaml:"state"`
	StateIn      []string `json:"state_in" yaml:"state_in"`
	Trend        string   `json:"trend" yaml:"trend"`
	Conclusion   string   `json:"conclusion" yaml:"conclusion"`
	ConfidenceGte float64 `json:"confidence_gte" yaml:"confidence_gte"`
	MinMatch     int     `json:"min_match" yaml:"min_match"`
	AnyOf        []Condition `json:"any_of" yaml:"any_of"`
}

// ActiveRule pairs a matching rule with its derived confidence
type ActiveRule struct {
	Rule      Rule
	Confidence float64
}

// DomainSnapshot is the output of domain evaluation
type DomainSnapshot struct {
	Domain          string  `json:"domain"`
	RuleID          string  `json:"rule_id"`
	Severity        string  `json:"severity"`
	Priority        int     `json:"priority"`
	Headline        string  `json:"headline"`
	RootCause       string  `json:"root_cause"`
	ServiceImpact   string  `json:"service_impact"`
	RuleConfidence  float64 `json:"rule_confidence"`
	EvaluatedAt     int64   `json:"evaluated_at"`
	// Anti-flap: when de-escalating, we hold until candidate is stable
	CandidateRuleID string `json:"candidate_rule_id,omitempty"`
	CandidateSince  int64  `json:"candidate_since,omitempty"`
}

// MetricEvaluationsMap maps metric pattern key (e.g. "cpu", "memory") to MetricEvaluation
type MetricEvaluationsMap map[string]MetricEvaluation

// EvaluateDomain runs the deterministic domain evaluator (rules engine)
func EvaluateDomain(domain string, metricEvals MetricEvaluationsMap, rules []Rule, previousSnapshot *DomainSnapshot, now int64) DomainSnapshot {
	activeRules := make([]ActiveRule, 0)

	for _, rule := range rules {
		if rule.Domain != domain {
			continue
		}
		if ruleMatches(rule, metricEvals) {
			conf := deriveRuleConfidence(rule, metricEvals)
			activeRules = append(activeRules, ActiveRule{Rule: rule, Confidence: conf})
		}
	}

	var candidate Rule
	if len(activeRules) == 0 {
		candidate = defaultHealthyRule(domain)
	} else {
		candidate = selectBest(activeRules)
	}

	finalRule, holdRuleID, holdSince := applyAntiflap(candidate, previousSnapshot, now)
	return buildDomainSnapshot(finalRule, metricEvals, now, holdRuleID, holdSince)
}

func ruleMatches(rule Rule, metricEvals MetricEvaluationsMap) bool {
	for _, cond := range rule.Conditions {
		if !conditionMatches(cond, metricEvals) {
			return false
		}
	}
	return true
}

func conditionMatches(cond Condition, metricEvals MetricEvaluationsMap) bool {
	if len(cond.AnyOf) > 0 {
		return anyOfMatches(cond.AnyOf, cond.MinMatch, metricEvals)
	}
	return simpleConditionMatches(cond, metricEvals)
}

func anyOfMatches(subs []Condition, minMatch int, metricEvals MetricEvaluationsMap) bool {
	count := 0
	for _, sub := range subs {
		if simpleConditionMatches(sub, metricEvals) {
			count++
		}
	}
	if minMatch <= 0 {
		minMatch = 1
	}
	return count >= minMatch
}

func simpleConditionMatches(cond Condition, metricEvals MetricEvaluationsMap) bool {
	eval, ok := metricEvals[cond.Metric]
	if !ok {
		return false
	}
	if !matchesState(cond, eval) || !matchesTrend(cond, eval) || !matchesConclusion(cond, eval) || !matchesConfidence(cond, eval) {
		return false
	}
	return true
}

func matchesState(cond Condition, eval MetricEvaluation) bool {
	if cond.State != "" && eval.State != cond.State {
		return false
	}
	for _, s := range cond.StateIn {
		if eval.State == s {
			return true
		}
	}
	return len(cond.StateIn) == 0
}

func matchesTrend(cond Condition, eval MetricEvaluation) bool {
	return cond.Trend == "" || eval.Trend == cond.Trend
}

func matchesConclusion(cond Condition, eval MetricEvaluation) bool {
	return cond.Conclusion == "" || eval.Conclusion == cond.Conclusion
}

func matchesConfidence(cond Condition, eval MetricEvaluation) bool {
	return cond.ConfidenceGte <= 0 || eval.Confidence >= cond.ConfidenceGte
}

func deriveRuleConfidence(rule Rule, metricEvals MetricEvaluationsMap) float64 {
	if len(rule.EvidenceMetrics) == 0 {
		return 0
	}
	var sum float64
	var n int
	for _, m := range rule.EvidenceMetrics {
		if e, ok := metricEvals[m]; ok {
			sum += e.Confidence
			n++
		}
	}
	if n == 0 {
		return 0
	}
	return sum / float64(n)
}

var severityOrder = map[string]int{
	"critical": 3,
	"degraded": 2,
	"healthy":  1,
}

func severityRank(s string) int {
	if r, ok := severityOrder[s]; ok {
		return r
	}
	return 0
}

// SeverityRank returns the numeric rank for severity comparison (higher = worse).
func SeverityRank(s string) int {
	return severityRank(s)
}

// ruleReferencesMetric returns true if the rule references the given metric in conditions or evidence.
func ruleReferencesMetric(rule Rule, metric string) bool {
	for _, m := range rule.EvidenceMetrics {
		if m == metric {
			return true
		}
	}
	for _, cond := range rule.Conditions {
		if conditionReferencesMetric(cond, metric) {
			return true
		}
	}
	return false
}

func conditionReferencesMetric(cond Condition, metric string) bool {
	if cond.Metric == metric {
		return true
	}
	for _, sub := range cond.AnyOf {
		if sub.Metric == metric {
			return true
		}
	}
	return false
}

// RulesForMetric returns rules that reference the given metric (valid for that metric type).
func RulesForMetric(rules []Rule, metric string) []Rule {
	out := make([]Rule, 0, len(rules))
	for _, r := range rules {
		if ruleReferencesMetric(r, metric) {
			out = append(out, r)
		}
	}
	return out
}

func selectBest(activeRules []ActiveRule) Rule {
	sort.Slice(activeRules, func(i, j int) bool {
		r1, r2 := activeRules[i], activeRules[j]
		if severityRank(r1.Rule.Severity) != severityRank(r2.Rule.Severity) {
			return severityRank(r1.Rule.Severity) > severityRank(r2.Rule.Severity)
		}
		if r1.Rule.Priority != r2.Rule.Priority {
			return r1.Rule.Priority > r2.Rule.Priority
		}
		return r1.Confidence > r2.Confidence
	})
	return activeRules[0].Rule
}

func defaultHealthyRule(domain string) Rule {
	return Rule{
		ID:         domain + ".healthy",
		Domain:     domain,
		Severity:   "healthy",
		Priority:   0,
		Headline:   "HEALTHY",
		RootCause:  "No issues detected",
		ServiceImpact: "System operating normally",
	}
}

// applyAntiflap returns (rule to use, holdRuleID for snapshot when in de-escalation hold, holdSince)
func applyAntiflap(candidate Rule, previous *DomainSnapshot, now int64) (Rule, string, int64) {
	if previous == nil {
		return candidate, "", 0
	}
	if candidate.ID == previous.RuleID {
		return candidate, previous.CandidateRuleID, previous.CandidateSince
	}
	if severityRank(candidate.Severity) > severityRank(previous.Severity) {
		return candidate, "", 0 // Escalate immediately
	}

	// De-escalation: require stability
	candidateSince := previous.CandidateSince
	candidateRuleID := previous.CandidateRuleID
	if candidateRuleID != candidate.ID {
		candidateSince = now
		candidateRuleID = candidate.ID
	}
	if now-candidateSince >= MinStableSec {
		return candidate, "", 0
	}
	// Keep previous; candidate not yet stable (pass through hold state for snapshot)
	return ruleFromSnapshot(previous), candidateRuleID, candidateSince
}

func ruleFromSnapshot(s *DomainSnapshot) Rule {
	return Rule{
		ID:             s.RuleID,
		Domain:         s.Domain,
		Severity:       s.Severity,
		Priority:       s.Priority,
		Headline:       s.Headline,
		RootCause:      s.RootCause,
		ServiceImpact:  s.ServiceImpact,
	}
}

func buildDomainSnapshot(rule Rule, metricEvals MetricEvaluationsMap, now int64, holdRuleID string, holdSince int64) DomainSnapshot {
	conf := 0.0
	if len(rule.EvidenceMetrics) > 0 {
		var sum float64
		n := 0
		for _, m := range rule.EvidenceMetrics {
			if e, ok := metricEvals[m]; ok {
				sum += e.Confidence
				n++
			}
		}
		if n > 0 {
			conf = sum / float64(n)
		}
	}

	snap := DomainSnapshot{
		Domain:         rule.Domain,
		RuleID:         rule.ID,
		Severity:       rule.Severity,
		Priority:       rule.Priority,
		Headline:       rule.Headline,
		RootCause:      rule.RootCause,
		ServiceImpact:  rule.ServiceImpact,
		RuleConfidence: conf,
		EvaluatedAt:    now,
	}

	if holdRuleID != "" && holdRuleID != rule.ID {
		snap.CandidateRuleID = holdRuleID
		snap.CandidateSince = holdSince
	}

	return snap
}
