/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Ctx, Query, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { Context } from "../context";
import { DefaultMarkupResDto } from "./types";

@Resolver()
export class GetDefaultMarkupResolver {
  @Query(() => DefaultMarkupResDto)
  @UseMiddleware(Authentication)
  async getDefaultMarkup(@Ctx() ctx: Context): Promise<DefaultMarkupResDto> {
    const { dataSources, baseURL } = ctx;
    return dataSources.dataSource.getDefaultMarkup(baseURL);
  }
}
