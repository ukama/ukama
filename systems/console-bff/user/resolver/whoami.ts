/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { WhoamiDto } from "./types";

@Resolver()
export class WhoamiResolver {
  @Query(() => WhoamiDto)
  async whoami(@Ctx() ctx: Context): Promise<WhoamiDto> {
    const { dataSources, headers } = ctx;
    return await dataSources.dataSource.whoami(headers.userId);
  }
}
