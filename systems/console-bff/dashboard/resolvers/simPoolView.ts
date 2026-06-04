/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import {
  Arg,
  Ctx,
  FieldResolver,
  Int,
  Query,
  Resolver,
  Root,
} from "type-graphql";

import { SIM_STATUS, SIM_TYPES } from "../../common/enums";
import type { AppContext } from "../../server/context";
import { ServiceUrlResolver } from "../baseUrls";
import { deriveSimPool } from "../derive";
import { runSection } from "../section";
import { PoolSimsSection, SimPoolStatsSection, SimPoolView } from "./types";

const MAX_POOL_SIMS = 100;

type SimPoolViewRoot = SimPoolView & { _urls: ServiceUrlResolver };

/**
 * SIM pool composite (plan §3.2). Serves both SIM pool screens (network
 * manage + business manage) with one query.
 */
@Resolver(() => SimPoolView)
export class SimPoolViewResolver {
  @Query(() => SimPoolView)
  simPoolView(
    @Arg("simType") simType: string,
    @Ctx() ctx: AppContext
  ): SimPoolView {
    return Object.assign(new SimPoolView(), {
      simType,
      _urls: new ServiceUrlResolver(ctx.headers.orgName),
    });
  }

  @FieldResolver(() => SimPoolStatsSection)
  async stats(
    @Root() root: SimPoolViewRoot,
    @Ctx() ctx: AppContext
  ): Promise<SimPoolStatsSection> {
    const { value, error } = await runSection("stats", async () => {
      const url = await root._urls.url("sim");
      const stats = await ctx.dataSources.sim.getSimPoolStats(
        url,
        root.simType
      );
      return { ...stats, ...deriveSimPool(stats) };
    });
    return { ...value, error };
  }

  @FieldResolver(() => PoolSimsSection)
  async sims(
    @Root() root: SimPoolViewRoot,
    @Ctx() ctx: AppContext,
    @Arg("limit", () => Int, { defaultValue: 20 }) limit: number
  ): Promise<PoolSimsSection> {
    const capped = Math.min(Math.max(limit, 1), MAX_POOL_SIMS);
    const { value, error } = await runSection("sims", async () => {
      const url = await root._urls.url("sim");
      const res = await ctx.dataSources.sim.getSimsFromPool(url, {
        type: root.simType as SIM_TYPES,
        status: SIM_STATUS.ALL,
      });
      return res.sims.slice(0, capped);
    });
    return { sims: value, error };
  }
}
