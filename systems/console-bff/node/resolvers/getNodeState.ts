/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { SUB_GRAPHS } from "../../common/configs";
import { openStore } from "../../common/storage";
import { getBaseURL } from "../../common/utils";
import { Context } from "../context";
import { NodeStateRes } from "./types";

@Resolver()
export class GetNodeStateResolver {
  @Query(() => NodeStateRes)
  async getNodeState(@Arg("id") id: string, @Ctx() context: Context) {
    const { dataSources, headers } = context;
    const store = openStore();
    const baseURL = await getBaseURL(
      SUB_GRAPHS.nodeState.name,
      headers.orgName,
      store
    );
    return await dataSources.dataSource.getNodeState(baseURL.message, id);
  }
}
