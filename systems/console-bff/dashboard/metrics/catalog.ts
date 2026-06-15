/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Metric catalog — the single source of presentation + mock-shape metadata
 * for every KPI/graph key the console renders. Each entry mirrors the fields
 * the metric service returns per key (unit/format/threshold/tick*), plus the
 * `label` the backend doesn't send and a mock seed (base/min/max/jitter/
 * trend). The console renders whatever the BFF returns — it owns none of this.
 *
 * Going live per metric: add its key to LIVE_METRIC_KEYS (then real values
 * flow; label/unit/threshold are still backfilled from here when the upstream
 * omits them).
 */

export interface MetricThresholdMeta {
  min: number;
  normal: number;
  max: number;
}

export interface MetricMeta {
  label: string;
  unit: string;
  /** Value formatting hint for the console: "number" | "decimal". */
  format: string;
  /** Mock seed: centre value and bounds. */
  base: number;
  min: number;
  max: number;
  jitter: number;
  trend: number;
  threshold?: MetricThresholdMeta;
}

const PCT: Pick<MetricMeta, "unit" | "format" | "min" | "max" | "threshold"> = {
  unit: "%",
  format: "number",
  min: 0,
  max: 100,
  threshold: { min: 0, normal: 80, max: 100 },
};
const TEMP: Pick<MetricMeta, "unit" | "format" | "min" | "max" | "threshold"> =
  {
    unit: "°C",
    format: "number",
    min: 0,
    max: 100,
    threshold: { min: 0, normal: 80, max: 100 },
  };
const DBM: Pick<MetricMeta, "unit" | "format" | "min" | "max" | "threshold"> = {
  unit: "dBm",
  format: "decimal",
  min: 0,
  max: 40,
  threshold: { min: 0, normal: 31, max: 34 },
};

export const METRIC_CATALOG: Record<string, MetricMeta> = {
  // --- node health ---
  uptime: {
    label: "Uptime",
    unit: "s",
    format: "number",
    base: 864000,
    min: 0,
    max: 2592000,
    jitter: 0.01,
    trend: 0.02,
  },
  cpu_temperature: {
    label: "Temp. (CPU)",
    base: 46,
    jitter: 0.12,
    trend: 0.05,
    ...TEMP,
  },
  fem1_temperature: {
    label: "FEM 1 temp.",
    base: 44,
    jitter: 0.12,
    trend: 0.04,
    ...TEMP,
  },
  fem2_temperature: {
    label: "FEM 2 temp.",
    base: 48,
    jitter: 0.12,
    trend: 0.04,
    ...TEMP,
  },
  memory: { label: "Memory", base: 52, jitter: 0.1, trend: 0.06, ...PCT },
  cpu: { label: "CPU", base: 38, jitter: 0.16, trend: 0.05, ...PCT },
  disk: {
    label: "Disk",
    unit: "MB",
    format: "number",
    base: 6200,
    min: 0,
    max: 16000,
    jitter: 0.04,
    trend: 0.08,
    threshold: { min: 0, normal: 12000, max: 16000 },
  },
  // --- customers ---
  subscribers_active: {
    label: "Active subscribers",
    unit: "",
    format: "number",
    base: 32,
    min: 0,
    max: 200,
    jitter: 0.18,
    trend: 0.12,
  },
  // --- network: cellular ---
  cellular_uplink: {
    label: "Cellular uplink",
    unit: "Mbps",
    format: "decimal",
    base: 18,
    min: 0,
    max: 50,
    jitter: 0.2,
    trend: 0.1,
    threshold: { min: 0, normal: 5, max: 30 },
  },
  cellular_downlink: {
    label: "Cellular downlink",
    unit: "Mbps",
    format: "decimal",
    base: 64,
    min: 0,
    max: 200,
    jitter: 0.2,
    trend: 0.12,
    threshold: { min: 0, normal: 60, max: 160 },
  },
  // --- network: backhaul ---
  backhaul_uplink: {
    label: "Backhaul uplink",
    unit: "Mbps",
    format: "decimal",
    base: 22,
    min: 0,
    max: 250,
    jitter: 0.18,
    trend: 0.1,
    threshold: { min: 0, normal: 10, max: 200 },
  },
  backhaul_downlink: {
    label: "Backhaul downlink",
    unit: "Mbps",
    format: "decimal",
    base: 70,
    min: 0,
    max: 250,
    jitter: 0.18,
    trend: 0.12,
    threshold: { min: 0, normal: 10, max: 200 },
  },
  backhaul_latency: {
    label: "Backhaul latency",
    unit: "ms",
    format: "decimal",
    base: 35,
    min: 0,
    max: 1050,
    jitter: 0.25,
    trend: 0.05,
    threshold: { min: 0, normal: 800, max: 1000 },
  },
  // --- site power / infrastructure ---
  site_uptime_percentage: {
    label: "Uptime",
    unit: "%",
    format: "number",
    base: 96,
    min: 80,
    max: 100,
    jitter: 0.03,
    trend: 0.01,
  },
  battery_charge: {
    label: "Available power",
    unit: "%",
    format: "number",
    base: 78,
    min: 0,
    max: 100,
    jitter: 0.08,
    trend: 0.05,
  },
  solar_panel_power: {
    label: "Solar power",
    unit: "W",
    format: "number",
    base: 320,
    min: 0,
    max: 600,
    jitter: 0.18,
    trend: 0.1,
  },
  controller_temperature: {
    label: "Controller temp.",
    unit: "°C",
    format: "number",
    base: 42,
    min: 0,
    max: 90,
    jitter: 0.12,
    trend: 0.05,
    threshold: { min: 0, normal: 60, max: 80 },
  },
  load_current: {
    label: "Load current",
    unit: "A",
    format: "decimal",
    base: 5,
    min: 0,
    max: 12,
    jitter: 0.16,
    trend: 0.06,
  },
  // --- radio ---
  power: { label: "TX power", base: 31, jitter: 0.06, trend: 0.02, ...DBM },
  pa_power: { label: "PA power", base: 30, jitter: 0.06, trend: 0.02, ...DBM },
  rx_power: { label: "RX power", base: 28, jitter: 0.06, trend: 0.02, ...DBM },
  tx_power: { label: "TX power", base: 31, jitter: 0.06, trend: 0.02, ...DBM },
};

const FALLBACK: MetricMeta = {
  label: "",
  unit: "",
  format: "number",
  base: 40,
  min: 0,
  max: 100,
  jitter: 0.12,
  trend: 0.04,
};

export const metricMeta = (key: string): MetricMeta =>
  METRIC_CATALOG[key] ?? { ...FALLBACK, label: key };

/**
 * Keys backed by a real metric endpoint (node-scoped: requested with a nodeId
 * so the gateway resolves the node type from the id). Each generic key below
 * exists under tnode/anode/cnode in default-metrics.yaml.
 *
 * Still mocked, pending more than a gate flip:
 *  - subscribers_active: node "customers" asks at tnode scope, but the key only
 *    exists under `system` (tnode has lte_active_ue / lte_subscribers).
 *  - site_uptime_percentage: no backing series in default-metrics.yaml.
 *  - battery_charge, solar_panel_power, controller_temperature, load_current:
 *    cnode metrics consumed by the Site screen without a nodeId (resolves to
 *    `system` → 404). Needs site→controller-node-id wiring first.
 */
export const LIVE_METRIC_KEYS = new Set<string>([
  // health / resources (all node types)
  "uptime",
  "cpu",
  "memory",
  "disk",
  "cpu_temperature",
  // tnode cellular
  "cellular_uplink",
  "cellular_downlink",
  // tnode / cnode backhaul
  "backhaul_uplink",
  "backhaul_downlink",
  "backhaul_latency",
  // tnode / cnode power
  "power",
  // anode radio
  "pa_power",
  "rx_power",
  "tx_power",
  // anode FEM health
  "fem1_temperature",
  "fem2_temperature",
]);

/** Mock unless explicitly disabled; never mock a key that has a live endpoint. */
const MOCK_ENABLED = process.env.MOCK_METRICS !== "false";
export const isMockKey = (key: string): boolean =>
  MOCK_ENABLED && !LIVE_METRIC_KEYS.has(key);
