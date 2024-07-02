/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { CBooleanResponse } from "../../common/types";
import { Context } from "../context";
import { NodeInput } from "./types";

@Resolver()
export class DetachNodeResolver {
  @Mutation(() => CBooleanResponse)
  async detachhNode(@Arg("data") data: NodeInput, @Ctx() context: Context) {
    const { dataSources, baseURL } = context;
    return await dataSources.dataSource.detachhNode(baseURL, {
      id: data.id,
    });
  }
}
