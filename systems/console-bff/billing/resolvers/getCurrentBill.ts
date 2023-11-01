/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { BillResponse } from "./types";

@Resolver()
export class GetCurrentBillResolver {
  @Query(() => BillResponse)
  async getCurrentBill(@Ctx() ctx: Context): Promise<BillResponse> {
    const { dataSources } = ctx;
    return dataSources.dataSource.getCurrentBill();
  }
}
