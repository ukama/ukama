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
  NodeLockStatus,
  nodeResourceKey,
  toNodeStatus,
  withBlocking,
} from "../logic";
import { NodeOperationStatusDto, NodeOperationStatusInputDto } from "./types";

/**
 * Lock status for a single node — used by the node detail page to render its
 * actions busy/available. `busy` reflects the node's own lock; `blocked` also
 * covers the other nodes of its site, which the backend checks before any
 * node action. Stateless: read-through calls per request. Every read fails
 * open (idle) so a lock-service blip never permanently disables the control.
 */
@Resolver()
export class GetNodeOperationStatusResolver {
  @Query(() => NodeOperationStatusDto)
  async getNodeOperationStatus(
    @Arg("data") data: NodeOperationStatusInputDto,
    @Ctx() ctx: AppContext
  ): Promise<NodeOperationStatusDto> {
    const opURL = await ctx.urls.url("operation");
    const readStatus = async (id: string): Promise<NodeLockStatus> => {
      try {
        const lock = await ctx.dataSources.operation.getResourceLock(
          opURL,
          nodeResourceKey(id)
        );
        return toNodeStatus({ id, lock });
      } catch {
        return toNodeStatus({ id, failed: true });
      }
    };

    const { nodeId, busy, operation } = await readStatus(data.nodeId);
    const self = { nodeId, busy, operation };
    if (busy) return withBlocking(self, []);

    const nodeURL = await ctx.urls.url("node");
    const siteId = await ctx.dataSources.node
      .getNode(nodeURL, { id: data.nodeId })
      .then(n => n.site?.siteId)
      .catch(() => undefined);
    if (!siteId) return withBlocking(self, []);

    const siteNodes = await ctx.dataSources.node
      .getNodes(nodeURL, { siteId })
      .then(res => res.nodes ?? [])
      .catch(() => []);
    const siblings = await Promise.all(
      siteNodes.filter(n => n.id !== data.nodeId).map(n => readStatus(n.id))
    );
    return withBlocking(self, siblings);
  }
}
