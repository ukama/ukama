/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { Nodes } from "./types";

@Resolver()
export class GetNodesForSiteResolver {
  @Query(() => Nodes)
  async getNodesForSite(
    @Arg("siteId") siteId: string,
    @Ctx() context: Context
  ) {
    const { dataSources, baseURL } = context;
    return await dataSources.dataSource.getNodesForSite(baseURL, siteId);
  }
}
