/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { VERSION } from "../../common/configs";
import { BaseRESTDataSource } from "../../common/datasource";
import {
  GetKpiTimeSeriesDto,
  GetKpiValuesDto,
  KpiTimeSeriesInput,
  KpiValueDto,
  KpiValuesInput,
  ScopeEntryDto,
} from "../resolvers/types/kpi";
import {
  GetPerformanceReportDto,
  PerformanceReportInput,
  ReportRowDto,
} from "../resolvers/types/report";
import { normalizeKpiValue, normalizeReportCell } from "../money";
import { mapAnalytics } from "./mapper";

const ANALYTICS = "analytics";

/**
 * Raw wire shape of one KPI value as rendered by the analytics gateway
 * (protojson, EmitUnpopulated → lowerCamelCase keys, `scope` as an object).
 */
interface RawKpiValue {
  kpi: string;
  value: number;
  span?: string;
  op?: string;
  from?: string;
  to?: string;
  type?: string;
  unit?: string;
  symbol?: string;
  isPartial?: boolean;
  scope?: Record<string, string>;
  trend?: KpiValueDto["trend"];
  computedAt?: string;
}

interface RawGetKpiValues {
  values?: RawKpiValue[];
}

/** Raw wire shape of one performance report row (`attributes` as an object). */
interface RawReportRow {
  entityId: string;
  attributes?: Record<string, string>;
  cells?: ReportRowDto["cells"];
  status?: string;
}

interface RawGetPerformanceReport {
  report: string;
  title?: string;
  span?: string;
  rows?: RawReportRow[];
}

/** Turn the protobuf `scope` map into a GraphQL-expressible entry list. */
const toScopeEntries = (scope?: Record<string, string>): ScopeEntryDto[] =>
  Object.entries(scope ?? {}).map(([key, value]) => ({ key, value }));

/**
 * Datasource for the analytics gateway's generic KPI read API. The gateway
 * exposes only self-describing KPI endpoints (no per-domain composite
 * endpoints); this BFF currently reads latest values via `/kpis/values`.
 */
class AnalyticsAPI extends BaseRESTDataSource {
  /**
   * GET /v1/analytics/kpis/values — latest rollup value per KPI/scope for a
   * span, with trend. `keys` is a CSV of KPI keys; `span`/`op` fall back to the
   * gateway's per-KPI defaults; `network_id` is an optional scope filter.
   */
  getKpiValues = async (
    baseURL: string,
    data: KpiValuesInput,
    currency?: string
  ): Promise<GetKpiValuesDto> => {
    this.baseURL = baseURL;

    const q = new URLSearchParams();
    q.append("keys", data.keys.join(","));
    if (data.span) q.append("span", data.span);
    if (data.op) q.append("op", data.op);
    if (data.networkId) q.append("network_id", data.networkId);

    const res = await this.callGet<RawGetKpiValues>(
      `kpis/values?${q.toString()}`
    );

    const values: KpiValueDto[] = (res.values ?? []).map(v =>
      normalizeKpiValue({ ...v, scope: toScopeEntries(v.scope) }, currency)
    );

    return { values };
  };

  /**
   * GET /v1/analytics/kpis/timeseries — one rollup value per span bucket over
   * the [from, to) window (e.g. weekly REVENUE for a 9-week trend). Same shape
   * as getKpiValues; each value carries its own from/to bucket boundaries.
   */
  getKpiTimeSeries = async (
    baseURL: string,
    data: KpiTimeSeriesInput,
    currency?: string
  ): Promise<GetKpiTimeSeriesDto> => {
    this.baseURL = baseURL;

    const q = new URLSearchParams();
    q.append("keys", data.keys.join(","));
    if (data.span) q.append("span", data.span);
    if (data.op) q.append("op", data.op);
    if (data.from) q.append("from", data.from);
    if (data.to) q.append("to", data.to);
    if (data.networkId) q.append("network_id", data.networkId);
    if (data.siteId) q.append("site_id", data.siteId);

    const res = await this.callGet<RawGetKpiValues>(
      `kpis/timeseries?${q.toString()}`
    );

    const values: KpiValueDto[] = (res.values ?? []).map(v =>
      normalizeKpiValue({ ...v, scope: toScopeEntries(v.scope) }, currency)
    );

    return { values };
  };

  /**
   * GET /v1/analytics/reports/{report} — resource performance table composed
   * from the latest available KPI values plus entity attributes. `span`/
   * `network_id` scope the table and `top` caps the number of rows.
   */
  getPerformanceReport = async (
    baseURL: string,
    data: PerformanceReportInput,
    currency?: string
  ): Promise<GetPerformanceReportDto> => {
    this.baseURL = baseURL;

    const q = new URLSearchParams();
    if (data.span) q.append("span", data.span);
    if (data.networkId) q.append("network_id", data.networkId);
    if (data.top != null) q.append("top", String(data.top));

    const qs = q.toString();
    const res = await this.callGet<RawGetPerformanceReport>(
      `reports/${encodeURIComponent(data.report)}${qs ? `?${qs}` : ""}`
    );

    const rows: ReportRowDto[] = (res.rows ?? []).map(r => ({
      entityId: r.entityId,
      attributes: toScopeEntries(r.attributes),
      cells: (r.cells ?? []).map(c => normalizeReportCell(c, currency)),
      status: r.status,
    }));

    return { report: res.report, title: res.title, span: res.span, rows };
  };

  /* ---------------- internal helpers ---------------- */

  private callGet = async <T>(path: string, label?: string): Promise<T> => {
    const url = `/${VERSION}/${ANALYTICS}/${path}`;
    this.logger.info(`Analytics ${label ?? path} [GET]: ${this.baseURL}${url}`);
    return this.get(url)
      .then(res => mapAnalytics<T>(res))
      .catch(error => {
        this.logger.error(
          `Error fetching analytics ${label ?? path}: ${error}`
        );
        throw error;
      });
  };
}

export default AnalyticsAPI;
