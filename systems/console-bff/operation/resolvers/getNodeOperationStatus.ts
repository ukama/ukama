/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Arg, Ctx, Query, Resolver } from "type-graphql";

import type { AppContext } from "../../server/context";
import { activeOperation, isLockBusy, nodeResourceKey } from "../logic";
import { NodeOperationStatusDto, NodeOperationStatusInputDto } from "./types";

/**
 * Lock status for a single node — used by the node detail page to render its
 * "Restart node" action busy/available. Stateless: one read-through to the
 * operation system per request. Fails open (idle) if the read errors so a
 * lock-service blip never permanently disables the control.
 */
@Resolver()
export class GetNodeOperationStatusResolver {
  @Query(() => NodeOperationStatusDto)
  async getNodeOperationStatus(
    @Arg("data") data: NodeOperationStatusInputDto,
    @Ctx() ctx: AppContext
  ): Promise<NodeOperationStatusDto> {
    const baseURL = await ctx.urls.url("operation");
    const lock = await ctx.dataSources.operation
      .getResourceLock(baseURL, nodeResourceKey(data.nodeId))
      .catch(() => undefined);

    return {
      nodeId: data.nodeId,
      busy: isLockBusy(lock),
      operation: activeOperation(lock),
    };
  }
}
