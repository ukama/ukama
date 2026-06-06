/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Ctx, Query, Resolver } from "type-graphql";

import type { AppContext } from "../../server/context";
import { Nodes } from "./types";

@Resolver()
export class GetNodesLocationResolver {
  @Query(() => Nodes)
  async getNodesLocation(@Ctx() context: AppContext) {
    const { dataSources } = context;
    const baseURL = await context.urls.url("node");
    const res = await dataSources.node.getNodes(baseURL, {});
    const nodes: Nodes = { nodes: [] };

    for (const node of res.nodes) {
      if (node.latitude && node.longitude) nodes.nodes.push({ ...node });
    }

    return nodes;
  }
}
