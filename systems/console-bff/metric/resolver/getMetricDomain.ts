/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { logger } from "../../common/logger";
import { Context } from "../context";
import { MetricDomain } from "./types";

@Resolver()
export class GetMetricDomainResolver {
  @Query(() => MetricDomain)
  async getMetricDomain(
    @Ctx() ctx: Context,
    @Arg("metricId") metricId: string,
    @Arg("nodeId") nodeId: string
  ): Promise<MetricDomain> {
    logger.info(
      `Getting metric domain for metric ${metricId} and node ${nodeId}`
    );
    const { dataSources, baseURL } = ctx;
    const metricDomain = await dataSources.dataSource.getMetricDomain(
      baseURL,
      metricId,
      nodeId
    );
    return metricDomain;
  }
}
