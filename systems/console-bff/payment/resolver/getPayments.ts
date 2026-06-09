/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Query, Resolver } from "type-graphql";

import type { AppContext } from "../../server/context";
import { GetPaymentsInputDto, PaymentsDto } from "./types";

@Resolver()
export class GetPaymentsResolver {
  @Query(() => PaymentsDto)
  async getPayments(
    @Arg("data") data: GetPaymentsInputDto,
    @Ctx() ctx: AppContext
  ): Promise<PaymentsDto> {
    const { dataSources } = ctx;
    const baseURL = await ctx.urls.url("payments");
    return dataSources.payment.getPayments(baseURL, data);
  }
}
