/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { SIM_TYPES } from "../../common/enums";
import { Context } from "../context";
import { SimsResDto } from "./types";

@Resolver()
export class GetSimsResolver {
  @Query(() => SimsResDto)
  async getSims(
    @Arg("type") type: SIM_TYPES,
    @Ctx() ctx: Context
  ): Promise<SimsResDto> {
    const { dataSources, baseURL } = ctx;
    return dataSources.dataSource.getSims(baseURL, type);
  }
}
