/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { SimDataUsage, SimDataUsages, SimUsagesInputDto } from "./types";

@Resolver()
export class GetDataUsagesResolver {
  @Query(() => SimDataUsages)
  async getDataUsages(
    @Arg("data") data: SimUsagesInputDto,
    @Ctx() ctx: Context
  ): Promise<SimDataUsages> {
    const { dataSources, baseURL } = ctx;
    const usages: SimDataUsage[] = [];
    for (const item of data.for) {
      usages.push(
        await dataSources.dataSource.getDataUsage(baseURL, {
          to: data.to,
          from: data.from,
          type: data.type,
          iccid: item.iccid,
          simId: item.simId,
        })
      );
    }
    return {
      usages: usages,
    };
  }
}
