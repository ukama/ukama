/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { OrgDto } from "./types";

@Resolver()
export class GetOrgResolver {
  @Query(() => OrgDto)
  async getOrg(@Ctx() ctx: Context): Promise<OrgDto> {
    const { dataSources, headers } = ctx;
    return dataSources.dataSource.getOrg(headers.orgName);
  }
}
