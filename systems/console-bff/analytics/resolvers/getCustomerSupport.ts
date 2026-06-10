/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Arg, Ctx, Query, Resolver } from "type-graphql";

import type { AppContext } from "../../server/context";
import { CustomerByIdInput, CustomerSupportDto } from "./customer.types";

@Resolver()
export class GetCustomerSupportResolver {
  @Query(() => CustomerSupportDto)
  async getCustomerSupport(
    @Arg("data") data: CustomerByIdInput,
    @Ctx() ctx: AppContext
  ): Promise<CustomerSupportDto> {
    const baseURL = await ctx.urls.url("analytics");
    return ctx.dataSources.analytics.getCustomerSupport(baseURL, data);
  }
}
