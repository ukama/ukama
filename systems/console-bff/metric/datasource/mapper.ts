/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { MetricAnalysis, MetricDomain } from "../resolver/types";

interface MetricAnalysisApiRes {
  aggregated: {
    computed_at: string;
    value: number;
    min: number;
    max: number;
    p95: number;
    mean: number;
    median: number;
    sample_count: number;
    aggregation: string;
    noise_estimate: number;
  };
  trend: { type: string; value: number };
  confidence: { value: number };
  projection: { type: string; eta_sec: number };
  state: string;
}

const mapToMetricAnalysis = (data: MetricAnalysisApiRes): MetricAnalysis => {
  return {
    aggregated: {
      computed_at: data.aggregated.computed_at,
      value: data.aggregated.value,
      min: data.aggregated.min,
      max: data.aggregated.max,
      p95: data.aggregated.p95,
      mean: data.aggregated.mean,
      median: data.aggregated.median,
      sample_count: data.aggregated.sample_count,
      aggregation: data.aggregated.aggregation,
      noise_estimate: data.aggregated.noise_estimate,
    },
    trend: {
      type: data.trend.type,
      value: data.trend.value,
    },
    confidence: {
      value: data.confidence.value,
    },
    projection: {
      type: data.projection.type,
      eta_sec: data.projection.eta_sec,
    },
    state: data.state,
  };
};

export const parseMetricAnalysisRes = (res: {
  data?: { result?: MetricAnalysisApiRes[] };
}): MetricAnalysis => {
  const data =
    res.data?.result?.[0] ?? (res as unknown as MetricAnalysisApiRes);
  return mapToMetricAnalysis(data);
};

interface MetricDomainApiRes {
  domain: {
    rule_id: string;
    severity: string;
    headline: string;
    root_cause: string;
    service_impact: string;
    rule_confidence: number;
    evaluated_at: string;
    computed_at: string;
  };
}

export const parseMetricDomainRes = (res: MetricDomainApiRes): MetricDomain => {
  const { domain } = res;
  return {
    rule_id: domain.rule_id,
    severity: domain.severity,
    headline: domain.headline,
    root_cause: domain.root_cause,
    service_impact: domain.service_impact,
    rule_confidence: domain.rule_confidence,
    evaluated_at: domain.evaluated_at,
    computed_at: domain.computed_at,
  };
};
