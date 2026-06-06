/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Query, Resolver } from "type-graphql";

import type { AppContext } from "../../server/context";
import { NodeStateRes } from "./types";

@Resolver()
export class GetNodeStateResolver {
  @Query(() => NodeStateRes)
  async getNodeState(@Arg("id") id: string, @Ctx() context: AppContext) {
    const { dataSources } = context;
    // Node state lives behind the "state" service (mapped to the node system).
    const baseURL = await context.urls.url("state");
    return await dataSources.node.getNodeState(baseURL, id);
  }
}
