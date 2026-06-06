/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Query, Resolver } from "type-graphql";

import type { AppContext } from "../../server/context";
import { Node, NodeInput } from "./types";

@Resolver()
export class GetNodeResolver {
  @Query(() => Node)
  async getNode(@Arg("data") data: NodeInput, @Ctx() context: AppContext) {
    const { dataSources } = context;
    const baseURL = await context.urls.url("node");
    return await dataSources.node.getNode(baseURL, { id: data.id });
  }
}
