/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Metric catalog — the single source of presentation metadata for every
 * KPI/graph key the console renders. Each entry carries the `label` the metric
 * service doesn't send, plus unit/format/threshold used to backfill a response
 * when the upstream omits them. Values themselves always come from the metric
 * service; this catalog never produces data. The console renders whatever the
 * BFF returns — it owns none of this.
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
  threshold?: MetricThresholdMeta;
}

const PCT: Pick<MetricMeta, "unit" | "format" | "threshold"> = {
  unit: "%",
  format: "number",
  threshold: { min: 0, normal: 80, max: 100 },
};
const TEMP: Pick<MetricMeta, "unit" | "format" | "threshold"> = {
  unit: "°C",
  format: "number",
  threshold: { min: 0, normal: 80, max: 100 },
};
const DBM: Pick<MetricMeta, "unit" | "format" | "threshold"> = {
  unit: "dBm",
  format: "decimal",
  threshold: { min: 0, normal: 31, max: 34 },
};

export const METRIC_CATALOG: Record<string, MetricMeta> = {
  // --- node health ---
  uptime: { label: "Uptime", unit: "s", format: "number" },
  cpu_temperature: { label: "Temp. (CPU)", ...TEMP },
  fem1_temperature: { label: "FEM 1 temp.", ...TEMP },
  fem2_temperature: { label: "FEM 2 temp.", ...TEMP },
  memory: { label: "Memory", ...PCT },
  cpu: { label: "CPU", ...PCT },
  disk: {
    label: "Disk",
    unit: "MB",
    format: "number",
    threshold: { min: 0, normal: 12000, max: 16000 },
  },
  // --- customers ---
  // tnode active subscribers — real series (trx_lte_core_active_ue).
  subscribers_active: {
    label: "Active customers",
    unit: "",
    format: "number",
    threshold: { min: 0, normal: 100, max: 1000 },
  },
  // --- network: cellular ---
  cellular_uplink: {
    label: "Cellular uplink",
    unit: "Mbps",
    format: "decimal",
    threshold: { min: 0, normal: 5, max: 30 },
  },
  cellular_downlink: {
    label: "Cellular downlink",
    unit: "Mbps",
    format: "decimal",
    threshold: { min: 0, normal: 60, max: 160 },
  },
  // --- network: backhaul ---
  backhaul_uplink: {
    label: "Backhaul uplink",
    unit: "Mbps",
    format: "decimal",
    threshold: { min: 0, normal: 10, max: 200 },
  },
  backhaul_downlink: {
    label: "Backhaul downlink",
    unit: "Mbps",
    format: "decimal",
    threshold: { min: 0, normal: 10, max: 200 },
  },
  backhaul_latency: {
    label: "Backhaul latency",
    unit: "ms",
    format: "decimal",
    threshold: { min: 0, normal: 800, max: 1000 },
  },
  // --- site power / infrastructure ---
  site_uptime_percentage: {
    label: "Uptime",
    unit: "%",
    format: "number",
  },
  battery_charge: {
    label: "Available power",
    unit: "%",
    format: "number",
  },
  solar_panel_power: {
    label: "Solar power",
    unit: "W",
    format: "number",
  },
  solar_panel_voltage: {
    label: "Solar voltage",
    unit: "V",
    format: "number",
    threshold: { min: 0, normal: 75, max: 100 },
  },
  solar_panel_current: {
    label: "Solar current",
    unit: "A",
    format: "number",
    threshold: { min: 0, normal: 5, max: 12 },
  },
  controller_temperature: {
    label: "Controller temp.",
    unit: "°C",
    format: "number",
    threshold: { min: 0, normal: 60, max: 80 },
  },
  load_current: {
    label: "Load current",
    unit: "A",
    format: "decimal",
  },
  // --- radio ---
  power: { label: "TX power", ...DBM },
  pa_power: { label: "PA power", ...DBM },
  rx_power: { label: "RX power", ...DBM },
  tx_power: { label: "TX power", ...DBM },
};

const FALLBACK: MetricMeta = {
  label: "",
  unit: "",
  format: "number",
};

export const metricMeta = (key: string): MetricMeta =>
  METRIC_CATALOG[key] ?? { ...FALLBACK, label: key };
