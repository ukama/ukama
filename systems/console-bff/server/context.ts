/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Unified per-request context for the consolidated API server
 * (CONSOLIDATION-DESIGN §4). Each domain gets a named datasource slot —
 * resolvers use `ctx.dataSources.<domain>` instead of the per-subgraph
 * `ctx.dataSources.dataSource`. The planning-tool resolvers additionally read
 * `ctx.prisma` when re-enabled (see README "Re-enabling planning-tool").
 */
import type { IncomingHttpHeaders } from "http";

import BillingAPI from "../billing/datasource/billing_api";
import { THeaders } from "../common/types";
import { parseExpressHeaders, parseToken } from "../common/utils";
import ComponentAPI from "../component/datasource/component_api";
import ControllerAPI from "../controller/datasource/controller_api";
import HealthAPI from "../health/datasource/health_api";
import InvitationAPI from "../invitation/datasource/invitation_api";
import MemberAPI from "../member/datasource/member_api";
import MetricAPI from "../metric/datasource/metric_api";
import NetworkAPI from "../network/datasource/network_api";
import NodeAPI from "../node/dataSource/node-api";
import NotificationAPI from "../notification/datasource/notification_api";
import OrgAPI from "../org/datasource/org_api";
import PackageAPI from "../package/datasource/package_api";
import PaymentAPI from "../payment/datasource/payment_api";
import RateAPI from "../rate/datasource/rate_api";
import ReportAPI from "../report/datasource/report_api";
import SimAPI from "../sim/datasource/sim_api";
import SiteAPI from "../site/datasource/site_api";
import SoftwareAPI from "../software/datasource/software_api";
import SubscriberAPI from "../subscriber/datasource/subscriber_api";
import UserAPI from "../user/datasource/user_api";

export interface AppDataSources {
  org: OrgAPI;
  user: UserAPI;
  network: NetworkAPI;
  site: SiteAPI;
  member: MemberAPI;
  invitation: InvitationAPI;
  node: NodeAPI;
  package: PackageAPI;
  rate: RateAPI;
  sim: SimAPI;
  subscriber: SubscriberAPI;
  controller: ControllerAPI;
  health: HealthAPI;
  software: SoftwareAPI;
  component: ComponentAPI;
  billing: BillingAPI;
  payment: PaymentAPI;
  report: ReportAPI;
  metric: MetricAPI;
  notification: NotificationAPI;
}

export interface AppContext {
  headers: THeaders;
  requestId: string;
  dataSources: AppDataSources;
  // NOTE: when planning-tool is re-enabled, add `prisma` here (see README
  // "Re-enabling planning-tool").
}

/**
 * Parses + verifies auth headers and decodes the signed token's claims so
 * resolvers get orgId/orgName/userId directly (the gateway and subgraphs
 * previously split this across two hops).
 */
export const buildHeaders = (reqHeaders: IncomingHttpHeaders): THeaders => {
  const headers = parseExpressHeaders(reqHeaders);
  if (headers.token) {
    headers.orgId = parseToken(headers.token, "orgId") ?? "";
    headers.userId = parseToken(headers.token, "userId") ?? "";
    headers.orgName = parseToken(headers.token, "orgName") ?? "";
  }
  return headers;
};

/** Per-request datasource instances (same lifecycle as the old subgraphs). */
export const buildDataSources = (): AppDataSources => ({
  org: new OrgAPI(),
  user: new UserAPI(),
  network: new NetworkAPI(),
  site: new SiteAPI(),
  member: new MemberAPI(),
  invitation: new InvitationAPI(),
  node: new NodeAPI(),
  package: new PackageAPI(),
  rate: new RateAPI(),
  sim: new SimAPI(),
  subscriber: new SubscriberAPI(),
  controller: new ControllerAPI(),
  health: new HealthAPI(),
  software: new SoftwareAPI(),
  component: new ComponentAPI(),
  billing: new BillingAPI(),
  payment: new PaymentAPI(),
  report: new ReportAPI(),
  metric: new MetricAPI(),
  notification: new NotificationAPI(),
});
