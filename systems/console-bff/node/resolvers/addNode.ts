/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { Context } from "../context";
import { AddNodeInput, Node } from "./types";

@Resolver()
export class AddNodeResolver {
  @Mutation(() => Node)
  async addNode(@Arg("data") data: AddNodeInput, @Ctx() context: Context) {
    const { dataSources } = context;
    return await dataSources.dataSource.addNode({
      id: data.id,
      name: data.name,
      orgId: data.orgId,
    });
  }
}
