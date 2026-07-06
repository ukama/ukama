/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Arg, Ctx, Query, Resolver } from "type-graphql";

import type { AppContext } from "../../server/context";
import {
  NodeLockRead,
  buildSiteActions,
  nodeResourceKey,
  toNodeStatus,
} from "../logic";
import { SiteOperationStatusDto, SiteOperationStatusInputDto } from "./types";

/**
 * Aggregated lock status for a site. A site is not itself lockable — it is a
 * set of nodes (tower/tnode, amplifier/anode, …), each locked independently
 * and released on its own async schedule. This resolver enumerates the site's
 * nodes and reads each node's lock in parallel, then derives site-level `busy`
 * and per-action availability so the console can render without knowing the
 * topology.
 *
 * Stateless (no background job): one fan-out of read-through calls per
 * request; the console owns the polling cadence. Per-node read failures fail
 * open (that node is treated as idle) and set `degraded` so the UI can show a
 * soft warning while correctness is still guarded by the backend's
 * conflict-on-start.
 */
@Resolver()
export class GetSiteOperationStatusResolver {
  @Query(() => SiteOperationStatusDto)
  async getSiteOperationStatus(
    @Arg("data") data: SiteOperationStatusInputDto,
    @Ctx() ctx: AppContext
  ): Promise<SiteOperationStatusDto> {
    const nodeURL = await ctx.urls.url("node");
    const opURL = await ctx.urls.url("operation");

    const nodesRes = await ctx.dataSources.node
      .getNodesForSite(nodeURL, data.siteId)
      .catch(() => ({ nodes: [] }));
    const nodes = nodesRes.nodes ?? [];

    const reads: NodeLockRead[] = await Promise.all(
      nodes.map(async n => {
        try {
          const lock = await ctx.dataSources.operation.getResourceLock(
            opURL,
            nodeResourceKey(n.id)
          );
          return { id: n.id, type: n.type, lock };
        } catch {
          return { id: n.id, type: n.type, failed: true };
        }
      })
    );

    const now = Date.now();
    const statuses = reads.map(r => toNodeStatus(r, now));

    return {
      siteId: data.siteId,
      busy: statuses.some(s => s.busy),
      degraded: reads.some(r => r.failed),
      nodes: statuses,
      actions: buildSiteActions(statuses),
    };
  }
}
