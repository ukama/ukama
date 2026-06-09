/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { logger } from "../../common/logger";
import { mapWithConcurrency } from "../../common/utils/concurrency";
import type { AppContext } from "../../server/context";
import { SimDataUsages, SimDto, SimUsagesInputDto } from "./types";

@Resolver()
export class GetDataUsagesResolver {
  @Query(() => SimDataUsages)
  async getDataUsages(
    @Arg("data") data: SimUsagesInputDto,
    @Ctx() ctx: AppContext
  ): Promise<SimDataUsages> {
    const { dataSources } = ctx;
    const baseURL = await ctx.urls.url("sim");

    const sims = await dataSources.sim.list(baseURL, {
      networkId: data.networkId,
      status: "active",
    });

    const simUsages: any =
      sims.sims
        .map((s: SimDto) => {
          if (s && s.id && s.package && s.package.id) {
            return {
              simId: s.id,
              iccid: s.iccid,
              packageEnd: s.package.endDate,
              packageStart: s.package.startDate,
            };
          }
          return null;
        })
        .filter(item => item !== null) ?? [];

    logger.info(`SimUsages: ${JSON.stringify(simUsages)}`);

    // Bounded fan-out: one upstream call per sim, max 10 in flight.
    const usages = await mapWithConcurrency(simUsages, (item: any) =>
      dataSources.sim.getDataUsage(baseURL, {
        type: data.type,
        iccid: item.iccid,
        simId: item.simId,
      })
    );

    return {
      usages,
    };
  }
}
