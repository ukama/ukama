/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { Node, NodeInput } from "./types";

@Resolver()
export class GetNodeResolver {
  @Query(() => Node)
  async getNode(@Arg("data") data: NodeInput, @Ctx() context: Context) {
    const { dataSources, baseURL } = context;
    return await dataSources.dataSource.getNode(baseURL, { id: data.id });
  }
}
