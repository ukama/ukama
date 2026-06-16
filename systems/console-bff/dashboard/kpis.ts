/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * KPI fan-out for composite sections (plan Phase 4): latest value per metric
 * key from the metric service, bounded concurrency, polled by the console —
 * no subscriptions/WS in v1 (BUILD-PLAN §5.1·4).
 */
import { logger } from "../common/logger";
import { mapWithConcurrency } from "../common/utils/concurrency";
import type MetricAPI from "../metric/datasource/metric_api";
import type { MetricThresholdMeta } from "./metrics/catalog";
import { isMockKey, metricMeta } from "./metrics/catalog";
import { mockLatest } from "./metrics/mock";

/** Metric keys per KPI section (see common/utils TYPE_KEYS_GROUPS). */
export const NETWORK_KPI_KEYS = [
  "network_uptime",
  "data_usage",
  "package_sales",
  "node_active_subscribers",
] as const;

export const SITE_KPI_KEYS = [
  "backhaul_latency",
  "backhaul_downlink",
  "controller_temperature",
  "load_current",
] as const;

export const SITE_POWER_KEYS = [
  "battery_charge",
  "solar_panel_voltage",
  "solar_panel_current",
  "solar_panel_power",
] as const;

export interface KpiEntry {
  key: string;
  value: number;
  timestamp: number;
  success: boolean;
  /** Presentation metadata so the console renders without a local catalog. */
  label?: string;
  unit?: string;
  format?: string;
  threshold?: MetricThresholdMeta | null;
}

export const fetchLatestKpis = async (
  metricApi: MetricAPI,
  baseURL: string,
  keys: readonly string[],
  opts: { nodeId?: string } = {}
): Promise<KpiEntry[]> => {
  const results = await mapWithConcurrency(keys, async key => {
    // One failing live key must not blank the whole section — degrade that
    // single KPI to success:false and let the rest render.
    try {
      if (isMockKey(key)) return { type: key, ...mockLatest(key, baseURL) };
      // Per-node KPIs must be node-scoped (the org-scoped /v1/metrics handler
      // hardcodes the `system` node type and 404s on node-only metrics).
      return opts.nodeId
        ? await metricApi.getNodeLatest(baseURL, key, opts.nodeId)
        : await metricApi.getLatestMetric(baseURL, key);
    } catch (e) {
      logger.warn(`[fetchLatestKpis] '${key}' failed: ${e}`);
      return {
        type: key,
        value: [Math.floor(Date.now() / 1000), 0] as [number, number],
        success: false,
      };
    }
  });
  return results.map(result => {
    const meta = metricMeta(result.type);
    return {
      key: result.type,
      timestamp: result.value[0],
      value: result.value[1],
      success: result.success,
      label: meta.label || result.type,
      unit: meta.unit,
      format: meta.format,
      threshold: meta.threshold ?? null,
    };
  });
};
