/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { InvoiceDto } from "./types";

@Resolver()
export class GetInvoiceByNetworkResolver {
  @Query(() => InvoiceDto)
  async getInvoicePDF(
    @Arg("invoiceId") invoiceId: string,
    @Ctx() ctx: Context
  ): Promise<any> {
    const { dataSources } = ctx;
    return dataSources.dataSource.GetInvoicePDF(invoiceId);
  }
}
