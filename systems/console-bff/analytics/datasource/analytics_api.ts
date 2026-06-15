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
  AnalyticsSiteInput,
  BillingSummaryDto,
  BusinessHomeDto,
  BusinessSiteDto,
  BusinessSitesDto,
  InventoryReadinessDto,
  PackagePerformanceDto,
  SalesOverviewDto,
} from "../resolvers/types/business";
import {
  RebuildRollupsInput,
  RebuildRollupsResultDto,
  RefreshInput,
  RefreshResultDto,
  RefreshStateDto,
} from "../resolvers/types/collector";
import {
  CustomerByIdInput,
  CustomerDetailDto,
  CustomerListDto,
  CustomerOverviewDto,
  CustomerSimsDto,
  CustomerSupportDto,
  SimPoolDto,
} from "../resolvers/types/customer";
import { HomeKpis, HomeLens, HomeViewInput } from "../resolvers/types/home";
import {
  AnalyticsNodeInput,
  MetricPanelDto,
  NetworkAlarmsDto,
  NetworkEventsDto,
  NetworkMetricsDto,
  NetworkNodeDto,
  NetworkNodesDto,
  NetworkOverviewDto,
  NetworkSiteDto,
  NetworkSitesDto,
  NetworkSupportSearchDto,
  NetworkTopologyDto,
  NodePoolDto,
} from "../resolvers/types/network";
import { AnalyticsWindowInput, KpiDto } from "../resolvers/types/shared";
import { mapAnalytics } from "./mapper";

const ANALYTICS = "analytics";

/**
 * Revenue KPI keys the console Revenue screen reads that the business
 * `GetSalesOverview` does not yet emit (it only emits `revenue`, `purchases`,
 * `avg_purchase`, `paid_customers`). See deriveRevenueKpis below.
 */
const REVENUE_SHIM_KEYS = [
  "revenue_collected",
  "revenue_month",
  "revenue_prev_month",
  "revenue_pending",
] as const;

const kpiValue = (
  o: SalesOverviewDto | undefined,
  key: string
): KpiDto | undefined => o?.kpis?.find(k => k.key === key);

/** Build a snake_case query string from a windowed/filtered/paginated input. */
const windowQuery = (data: AnalyticsWindowInput): string => {
  const q = new URLSearchParams();
  if (data.networkId) q.append("network_id", data.networkId);
  if (data.siteId) q.append("site_id", data.siteId);
  if (data.status) q.append("status", data.status);
  if (data.query) q.append("q", data.query);
  if (data.period) q.append("period", data.period);
  if (data.from) q.append("from", data.from);
  if (data.to) q.append("to", data.to);
  if (data.timezone) q.append("timezone", data.timezone);
  if (data.page) q.append("page", String(data.page));
  if (data.pageSize) q.append("page_size", String(data.pageSize));
  return q.toString();
};

/** Window-only query (period/from/to/timezone) for by-id detail endpoints. */
const detailQuery = (data: {
  period?: string;
  from?: string;
  to?: string;
  timezone?: string;
}): string => {
  const q = new URLSearchParams();
  if (data.period) q.append("period", data.period);
  if (data.from) q.append("from", data.from);
  if (data.to) q.append("to", data.to);
  if (data.timezone) q.append("timezone", data.timezone);
  return q.toString();
};

class AnalyticsAPI extends BaseRESTDataSource {
  /* ---------------- Home (shared business / network) ---------------- */

  // KPI strip for either home screen. Routes by lens to the lens's overview
  // endpoint and returns just its `kpis` (both lenses emit [KpiDto]).
  getHomeKpis = async (
    baseURL: string,
    data: HomeViewInput
  ): Promise<HomeKpis> => {
    this.baseURL = baseURL;
    const path =
      data.lens === HomeLens.NETWORK ? "network/overview" : "business/home";
    const res = await this.callGet<{ kpis?: KpiDto[] }>(
      `${path}?${windowQuery(data)}`
    );
    return { kpis: res.kpis ?? [] };
  };

  /* ---------------- Business ---------------- */

  getBusinessHome = async (
    baseURL: string,
    data: AnalyticsWindowInput
  ): Promise<BusinessHomeDto> => {
    this.baseURL = baseURL;
    return this.callGet<BusinessHomeDto>(`business/home?${windowQuery(data)}`);
  };

  getSalesOverview = async (
    baseURL: string,
    data: AnalyticsWindowInput
  ): Promise<SalesOverviewDto> => {
    this.baseURL = baseURL;
    const overview = await this.callGet<SalesOverviewDto>(
      `business/sales/overview?${windowQuery(data)}`
    );
    // TODO(analytics-backend): drop once the backend emits the revenue_* KPIs.
    return this.deriveRevenueKpis(overview, data);
  };

  /**
   * TODO(analytics-backend): the business `GetSalesOverview` only emits
   * `revenue`/`purchases`/`avg_purchase`/`paid_customers`, but the console
   * Revenue screen also needs `revenue_collected`, `revenue_month`,
   * `revenue_prev_month` and `revenue_pending`. Until the business service emits
   * these, we derive them from REAL payment data by re-querying the overview
   * endpoint for the relevant windows (this month / last month / all-time) and
   * reading each window's `revenue` KPI — so the values always match the data in
   * the analytics DB. `revenue_pending` has no backend source yet, so it stays 0.
   * Remove this shim once the backend provides these keys.
   * Tracking: systems/console-bff/docs/analytics-backend-gaps.md
   */
  private deriveRevenueKpis = async (
    overview: SalesOverviewDto,
    data: AnalyticsWindowInput
  ): Promise<SalesOverviewDto> => {
    const have = new Set((overview.kpis ?? []).map(k => k.key));
    if (REVENUE_SHIM_KEYS.every(k => have.has(k))) return overview;

    const now = new Date();
    const startOfMonth = new Date(
      Date.UTC(now.getUTCFullYear(), now.getUTCMonth(), 1)
    ).toISOString();
    const startOfPrevMonth = new Date(
      Date.UTC(now.getUTCFullYear(), now.getUTCMonth() - 1, 1)
    ).toISOString();

    const q = (extra: Partial<AnalyticsWindowInput>): string =>
      `business/sales/overview?${windowQuery({ networkId: data.networkId, ...extra })}`;

    // Each window's `revenue` KPI is computed by the backend from real payment
    // events, so these reflect exactly what's in the analytics DB.
    const safe = async (
      path: string
    ): Promise<SalesOverviewDto | undefined> => {
      try {
        return await this.callGet<SalesOverviewDto>(path);
      } catch {
        return undefined;
      }
    };

    const [monthRes, prevRes, allRes] = await Promise.all([
      safe(q({ period: "month" })),
      safe(q({ period: "custom", from: startOfPrevMonth, to: startOfMonth })),
      safe(
        q({
          period: "custom",
          from: "2000-01-01T00:00:00Z",
          to: now.toISOString(),
        })
      ),
    ]);

    const monthRev = kpiValue(monthRes, "revenue");
    const prevRev = kpiValue(prevRes, "revenue");
    const allRev = kpiValue(allRes, "revenue");

    const fmt = (v: number): string =>
      `$${(Math.round(v * 100) / 100).toLocaleString()}`;
    const kpi = (
      key: string,
      value: number,
      delta?: number,
      stale = false
    ): KpiDto => ({
      key,
      value: Math.round(value * 100) / 100,
      formatted: fmt(value),
      delta,
      deltaPeriod: delta != null ? "month" : undefined,
      stale,
    });

    const derived: Record<string, KpiDto> = {
      revenue_collected: kpi("revenue_collected", allRev?.value ?? 0),
      revenue_month: kpi(
        "revenue_month",
        monthRev?.value ?? 0,
        monthRev?.delta
      ),
      revenue_prev_month: kpi("revenue_prev_month", prevRev?.value ?? 0),
      // No pending-payment source in analytics yet → 0, flagged as mock.
      revenue_pending: kpi("revenue_pending", 0, undefined, true),
    };

    const merged = [...(overview.kpis ?? [])];
    for (const key of REVENUE_SHIM_KEYS) {
      if (!have.has(key)) merged.push(derived[key]);
    }

    return { ...overview, kpis: merged };
  };

  getPackagePerformance = async (
    baseURL: string,
    data: AnalyticsWindowInput
  ): Promise<PackagePerformanceDto> => {
    this.baseURL = baseURL;
    return this.callGet<PackagePerformanceDto>(
      `business/sales/packages?${windowQuery(data)}`
    );
  };

  getBillingSummary = async (
    baseURL: string,
    data: AnalyticsWindowInput
  ): Promise<BillingSummaryDto> => {
    this.baseURL = baseURL;
    return this.callGet<BillingSummaryDto>(
      `business/billing?${windowQuery(data)}`
    );
  };

  getBusinessSites = async (
    baseURL: string,
    data: AnalyticsWindowInput
  ): Promise<BusinessSitesDto> => {
    this.baseURL = baseURL;
    return this.callGet<BusinessSitesDto>(
      `business/sites?${windowQuery(data)}`
    );
  };

  getBusinessSite = async (
    baseURL: string,
    data: AnalyticsSiteInput
  ): Promise<BusinessSiteDto> => {
    this.baseURL = baseURL;
    return this.callGet<BusinessSiteDto>(
      `business/sites/${data.siteId}?${detailQuery(data)}`
    );
  };

  getInventoryReadiness = async (
    baseURL: string,
    data: AnalyticsWindowInput
  ): Promise<InventoryReadinessDto> => {
    this.baseURL = baseURL;
    return this.callGet<InventoryReadinessDto>(
      `business/inventory?${windowQuery(data)}`
    );
  };

  /* ---------------- Customers ---------------- */

  getCustomerOverview = async (
    baseURL: string,
    data: AnalyticsWindowInput
  ): Promise<CustomerOverviewDto> => {
    this.baseURL = baseURL;
    return this.callGet<CustomerOverviewDto>(
      `customers/overview?${windowQuery(data)}`
    );
  };

  listCustomers = async (
    baseURL: string,
    data: AnalyticsWindowInput
  ): Promise<CustomerListDto> => {
    this.baseURL = baseURL;
    return this.callGet<CustomerListDto>(`customers/list?${windowQuery(data)}`);
  };

  searchCustomers = async (
    baseURL: string,
    data: AnalyticsWindowInput
  ): Promise<CustomerListDto> => {
    this.baseURL = baseURL;
    return this.callGet<CustomerListDto>(
      `customers/search?${windowQuery(data)}`
    );
  };

  getSims = async (
    baseURL: string,
    data: AnalyticsWindowInput
  ): Promise<CustomerSimsDto> => {
    this.baseURL = baseURL;
    return this.callGet<CustomerSimsDto>(`customers/sims?${windowQuery(data)}`);
  };

  getSimPool = async (
    baseURL: string,
    data: AnalyticsWindowInput
  ): Promise<SimPoolDto> => {
    this.baseURL = baseURL;
    return this.callGet<SimPoolDto>(`customers/sim-pool?${windowQuery(data)}`);
  };

  getCustomer = async (
    baseURL: string,
    data: CustomerByIdInput
  ): Promise<CustomerDetailDto> => {
    this.baseURL = baseURL;
    return this.callGet<CustomerDetailDto>(
      `customers/${data.customerId}?${detailQuery(data)}`
    );
  };

  getCustomerSupport = async (
    baseURL: string,
    data: CustomerByIdInput
  ): Promise<CustomerSupportDto> => {
    this.baseURL = baseURL;
    return this.callGet<CustomerSupportDto>(
      `customers/${data.customerId}/support`
    );
  };

  /* ---------------- Network ---------------- */

  getNetworkOverview = async (
    baseURL: string,
    data: AnalyticsWindowInput
  ): Promise<NetworkOverviewDto> => {
    this.baseURL = baseURL;
    return this.callGet<NetworkOverviewDto>(
      `network/overview?${windowQuery(data)}`
    );
  };

  getTopology = async (
    baseURL: string,
    data: AnalyticsWindowInput
  ): Promise<NetworkTopologyDto> => {
    this.baseURL = baseURL;
    return this.callGet<NetworkTopologyDto>(
      `network/topology?${windowQuery(data)}`
    );
  };

  getNetworkSites = async (
    baseURL: string,
    data: AnalyticsWindowInput
  ): Promise<NetworkSitesDto> => {
    this.baseURL = baseURL;
    return this.callGet<NetworkSitesDto>(`network/sites?${windowQuery(data)}`);
  };

  getNetworkSite = async (
    baseURL: string,
    data: AnalyticsSiteInput
  ): Promise<NetworkSiteDto> => {
    this.baseURL = baseURL;
    return this.callGet<NetworkSiteDto>(
      `network/sites/${data.siteId}?${detailQuery(data)}`
    );
  };

  getNetworkNodes = async (
    baseURL: string,
    data: AnalyticsWindowInput
  ): Promise<NetworkNodesDto> => {
    this.baseURL = baseURL;
    return this.callGet<NetworkNodesDto>(`network/nodes?${windowQuery(data)}`);
  };

  getNetworkNode = async (
    baseURL: string,
    data: AnalyticsNodeInput
  ): Promise<NetworkNodeDto> => {
    this.baseURL = baseURL;
    return this.callGet<NetworkNodeDto>(
      `network/nodes/${data.nodeId}?${detailQuery(data)}`
    );
  };

  getNodePool = async (
    baseURL: string,
    data: AnalyticsWindowInput
  ): Promise<NodePoolDto> => {
    this.baseURL = baseURL;
    return this.callGet<NodePoolDto>(`network/node-pool?${windowQuery(data)}`);
  };

  getRadio = async (
    baseURL: string,
    data: AnalyticsWindowInput
  ): Promise<MetricPanelDto> => {
    this.baseURL = baseURL;
    return this.callGet<MetricPanelDto>(`network/radio?${windowQuery(data)}`);
  };

  getBackhaul = async (
    baseURL: string,
    data: AnalyticsWindowInput
  ): Promise<MetricPanelDto> => {
    this.baseURL = baseURL;
    return this.callGet<MetricPanelDto>(
      `network/backhaul?${windowQuery(data)}`
    );
  };

  getPower = async (
    baseURL: string,
    data: AnalyticsWindowInput
  ): Promise<MetricPanelDto> => {
    this.baseURL = baseURL;
    return this.callGet<MetricPanelDto>(`network/power?${windowQuery(data)}`);
  };

  getAlarms = async (
    baseURL: string,
    data: AnalyticsWindowInput
  ): Promise<NetworkAlarmsDto> => {
    this.baseURL = baseURL;
    return this.callGet<NetworkAlarmsDto>(
      `network/alarms?${windowQuery(data)}`
    );
  };

  getMetrics = async (
    baseURL: string,
    data: AnalyticsWindowInput
  ): Promise<NetworkMetricsDto> => {
    this.baseURL = baseURL;
    return this.callGet<NetworkMetricsDto>(
      `network/metrics?${windowQuery(data)}`
    );
  };

  getEvents = async (
    baseURL: string,
    data: AnalyticsWindowInput
  ): Promise<NetworkEventsDto> => {
    this.baseURL = baseURL;
    return this.callGet<NetworkEventsDto>(
      `network/events?${windowQuery(data)}`
    );
  };

  supportSearch = async (
    baseURL: string,
    data: AnalyticsWindowInput
  ): Promise<NetworkSupportSearchDto> => {
    this.baseURL = baseURL;
    return this.callGet<NetworkSupportSearchDto>(
      `network/support/search?${windowQuery(data)}`
    );
  };

  /* ---------------- Collector ---------------- */

  refreshAnalytics = async (
    baseURL: string,
    data: RefreshInput
  ): Promise<RefreshResultDto> => {
    this.baseURL = baseURL;
    return this.callPost<RefreshResultDto>("collector/refresh", {
      source: data.source,
    });
  };

  getRefreshState = async (baseURL: string): Promise<RefreshStateDto> => {
    this.baseURL = baseURL;
    return this.callGet<RefreshStateDto>("collector/state");
  };

  rebuildRollups = async (
    baseURL: string,
    data: RebuildRollupsInput
  ): Promise<RebuildRollupsResultDto> => {
    this.baseURL = baseURL;
    return this.callPost<RebuildRollupsResultDto>("collector/rollups/rebuild", {
      family: data.family,
      from: data.from,
      to: data.to,
    });
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

  private callPost = async <T>(
    path: string,
    body: Record<string, unknown>
  ): Promise<T> => {
    const url = `/${VERSION}/${ANALYTICS}/${path}`;
    this.logger.info(`Analytics ${path} [POST]: ${this.baseURL}${url}`);
    return this.post(url, { body })
      .then(res => mapAnalytics<T>(res))
      .catch(error => {
        this.logger.error(`Error posting analytics ${path}: ${error}`);
        throw error;
      });
  };
}

export default AnalyticsAPI;
