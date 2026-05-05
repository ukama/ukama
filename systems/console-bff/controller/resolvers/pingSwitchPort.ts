/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { CBooleanResponse } from "../../common/types";
import { Context } from "../context";
import { PingSwitchPortInputDto } from "./types";

@Resolver()
export class PingSwitchPortResolver {
  @Query(() => String, { nullable: true })
  async pingSwitchPort(
    @Arg("data") data: PingSwitchPortInputDto,
    @Ctx() ctx: Context
  ): Promise<CBooleanResponse> {
    const { dataSources, baseURL } = ctx;
    return dataSources.dataSource.pingSwitchPort(baseURL, data);
  }
}
