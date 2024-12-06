/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { GetReportResDto } from "./types";

@Resolver()
export class GetReportResolver {
  @Query(() => GetReportResDto)
  async getReport(
    @Arg("id") id: string,
    @Ctx() ctx: Context
  ): Promise<GetReportResDto> {
    const { dataSources, baseURL } = ctx;
    return dataSources.dataSource.getReport(baseURL, id);
  }
}
