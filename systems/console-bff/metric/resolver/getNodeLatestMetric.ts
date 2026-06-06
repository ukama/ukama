/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Query, Resolver } from "type-graphql";

import type { AppContext } from "../../server/context";
import { GetNodeLatestMetricInput, NodeLatestMetric } from "./types";

@Resolver()
export class GetNodeLatestMetricResolver {
  @Query(() => NodeLatestMetric)
  async getNodeLatestMetric(
    @Arg("data") data: GetNodeLatestMetricInput,
    @Ctx() ctx: AppContext
  ): Promise<NodeLatestMetric> {
    const { dataSources } = ctx;
    const baseURL = await ctx.urls.url("metrics");
    return await dataSources.metric.getNodeLatestMetric(baseURL, data);
  }
}
