/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { Context } from "../context";
import { AddSiteInputDto, SiteDto } from "./types";

@Resolver()
export class AddSiteResolver {
  @Mutation(() => SiteDto)
  async addSite(
    @Arg("data") data: AddSiteInputDto,
    @Ctx() ctx: Context
  ): Promise<SiteDto> {
    const { dataSources, baseURL } = ctx;
    return dataSources.dataSource.addSite(baseURL, data);
  }
}
