/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { SIM_STATUS } from "../../common/enums";
import { Context } from "../context";
import { GetSimsInput, SimsResDto } from "./types";

@Resolver()
export class GetSimsResolver {
  @Query(() => SimsResDto)
  async getSims(
    @Arg("data") data: GetSimsInput,
    @Ctx() ctx: Context
  ): Promise<SimsResDto> {
    const { dataSources, baseURL } = ctx;
    const sims = await dataSources.dataSource.getSims(baseURL, data.type);
    if (data.status === SIM_STATUS.ASSIGNED) {
      return { sim: sims.sim.filter(sim => sim.isAllocated === true) };
    } else if (data.status === SIM_STATUS.UNASSIGNED) {
      return { sim: sims.sim.filter(sim => sim.isAllocated === false) };
    } else {
      return sims;
    }
  }
}
