/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Arg, Ctx, Query, Resolver } from "type-graphql";

import type { AppContext } from "../../server/context";
import { GetNodeInterfacesInputDto, NodeInterfaces } from "./types";

@Resolver()
export class GetNodeInterfaces {
  @Query(() => NodeInterfaces)
  async getNodeInterfaces(
    @Ctx() ctx: AppContext,
    @Arg("data") data: GetNodeInterfacesInputDto
  ): Promise<NodeInterfaces> {
    // health resolves to the node gateway (nodeGwIp:nodeGwPort) via isForNodeGw.
    const baseURL = await ctx.urls.url("health");
    return ctx.dataSources.health.getInterfaces(baseURL, data);
  }
}
