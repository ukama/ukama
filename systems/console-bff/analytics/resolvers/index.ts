/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { NonEmptyArray } from "type-graphql";

import { GetAlarmsResolver } from "./getAlarms";
import { GetBackhaulResolver } from "./getBackhaul";
import { GetBillingSummaryResolver } from "./getBillingSummary";
import { GetBusinessHomeResolver } from "./getBusinessHome";
import { GetBusinessSiteResolver } from "./getBusinessSite";
import { GetBusinessSitesResolver } from "./getBusinessSites";
import { GetCustomerResolver } from "./getCustomer";
import { GetCustomerOverviewResolver } from "./getCustomerOverview";
import { GetCustomerSimsResolver } from "./getCustomerSims";
import { GetCustomerSupportResolver } from "./getCustomerSupport";
import { GetEventsResolver } from "./getEvents";
import { GetHomeKpisResolver } from "./getHomeKpis";
import { GetInventoryReadinessResolver } from "./getInventoryReadiness";
import { GetMetricsResolver } from "./getMetrics";
import { GetNetworkNodeResolver } from "./getNetworkNode";
import { GetNetworkNodesResolver } from "./getNetworkNodes";
import { GetNetworkOverviewResolver } from "./getNetworkOverview";
import { GetNetworkSiteResolver } from "./getNetworkSite";
import { GetNetworkSitesResolver } from "./getNetworkSites";
import { GetNodePoolResolver } from "./getNodePool";
import { GetPackagePerformanceResolver } from "./getPackagePerformance";
import { GetPowerResolver } from "./getPower";
import { GetRadioResolver } from "./getRadio";
import { GetRefreshStateResolver } from "./getRefreshState";
import { GetSalesOverviewResolver } from "./getSalesOverview";
import { GetSimPoolResolver } from "./getSimPool";
import { GetTopologyResolver } from "./getTopology";
import { ListCustomersResolver } from "./listCustomers";
import { RebuildRollupsResolver } from "./rebuildRollups";
import { RefreshAnalyticsResolver } from "./refreshAnalytics";
import { SearchCustomersResolver } from "./searchCustomers";
import { SupportSearchResolver } from "./supportSearch";

const resolvers: NonEmptyArray<any> = [
  // Home (shared business / network)
  GetHomeKpisResolver,
  // Business backend
  GetBusinessHomeResolver,
  GetSalesOverviewResolver,
  GetPackagePerformanceResolver,
  GetBillingSummaryResolver,
  GetBusinessSitesResolver,
  GetBusinessSiteResolver,
  GetInventoryReadinessResolver,
  // Customers backend
  GetCustomerOverviewResolver,
  ListCustomersResolver,
  SearchCustomersResolver,
  GetCustomerSimsResolver,
  GetSimPoolResolver,
  GetCustomerResolver,
  GetCustomerSupportResolver,
  // Network backend
  GetNetworkOverviewResolver,
  GetTopologyResolver,
  GetNetworkSitesResolver,
  GetNetworkSiteResolver,
  GetNetworkNodesResolver,
  GetNetworkNodeResolver,
  GetNodePoolResolver,
  GetRadioResolver,
  GetBackhaulResolver,
  GetPowerResolver,
  GetAlarmsResolver,
  GetMetricsResolver,
  GetEventsResolver,
  SupportSearchResolver,
  // Collector backend
  RefreshAnalyticsResolver,
  GetRefreshStateResolver,
  RebuildRollupsResolver,
];

export default resolvers;
