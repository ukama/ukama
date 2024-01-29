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
export class GetInvoiceResolver {
  @Query(() => InvoiceDto)
  async getInvoice(
    @Arg("invoiceId") invoiceId: string,
    @Ctx() ctx: Context
  ): Promise<InvoiceDto> {
    const { dataSources } = ctx;
    return dataSources.dataSource.getInvoice(invoiceId);
  }
}
