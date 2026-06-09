/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import type { AppContext } from "../../server/context";
import { Node, UpdateNodeStateInput } from "./types";

@Resolver()
export class UpdateNodeStateResolver {
  @Mutation(() => Node)
  async updateNodeState(
    @Arg("data") data: UpdateNodeStateInput,
    @Ctx() context: AppContext
  ) {
    const { dataSources } = context;
    const baseURL = await context.urls.url("node");
    return await dataSources.node.updateNodeState(baseURL, {
      id: data.id,
      state: data.state,
    });
  }
}
