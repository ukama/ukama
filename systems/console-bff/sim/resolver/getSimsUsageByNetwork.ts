/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { mapWithConcurrency } from "../../common/utils/concurrency";
import type { AppContext } from "../../server/context";
import { SimDto, SimUsageItem } from "./types";

/** cdr_type for data usage (systems/common/ukama CdrTypeData). */
const CDR_TYPE_DATA = "data";

/**
 * Data usage for every SIM on a network, as a list of <simId>:<usage> pairs.
 *
 * Abstracts the N-per-SIM fan-out away from the console: it lists the
 * network's SIMs once, then fetches usage per SIM with bounded concurrency
 * (max 10 in flight) and returns a flat map the console can index by SIM id.
 * A per-SIM usage read that fails yields "0" rather than dropping the SIM.
 */
@Resolver()
export class GetSimsUsageByNetworkResolver {
  @Query(() => [SimUsageItem])
  async getSimsUsageByNetwork(
    @Arg("networkId") networkId: string,
    @Ctx() ctx: AppContext
  ): Promise<SimUsageItem[]> {
    const baseURL = await ctx.urls.url("sim");

    const sims = await ctx.dataSources.sim.list(baseURL, {
      networkId,
      status: "service_on",
    });
    const items = (sims.sims ?? []).filter((s: SimDto) => s?.id && s?.iccid);

    return mapWithConcurrency(items, (s: SimDto) =>
      ctx.dataSources.sim
        .getDataUsage(baseURL, {
          simId: s.id,
          iccid: s.iccid,
          type: CDR_TYPE_DATA,
        })
        .then(u => ({ simId: s.id, usage: u.usage }))
        .catch(() => ({ simId: s.id, usage: "0" }))
    );
  }
}
