/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import type { AppContext } from "../../server/context";
import { AddNodeInput, Node } from "./types";

@Resolver()
export class AddNodeResolver {
  @Mutation(() => Node)
  async addNode(@Arg("data") data: AddNodeInput, @Ctx() context: AppContext) {
    const { dataSources } = context;
    const baseURL = await context.urls.url("node");
    return await dataSources.node.addNode(baseURL, {
      id: data.id,
      name: data.name,
    });
  }
}
