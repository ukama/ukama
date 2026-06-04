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
 * `ctx.dataSources.dataSource`. Slots are added as Phase B migrates each
 * module batch.
 */
import type { IncomingHttpHeaders } from "http";

import { THeaders } from "../common/types";
import { parseExpressHeaders, parseToken } from "../common/utils";
import OrgAPI from "../org/datasource/org_api";

export interface AppContext {
  headers: THeaders;
  requestId: string;
  dataSources: {
    org: OrgAPI;
    // …one key per module, extended per Phase B batch (user, network, …)
  };
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
export const buildDataSources = (): AppContext["dataSources"] => ({
  org: new OrgAPI(),
});
