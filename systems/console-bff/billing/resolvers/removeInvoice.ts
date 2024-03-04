/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { CBooleanResponse } from "../../common/types";
import { Context } from "../context";

@Resolver()
export class RemoveInvoiceResolver {
  @Mutation(() => CBooleanResponse)
  async removeInvoice(
    @Arg("invoiceId") invoiceId: string,
    @Ctx() ctx: Context
  ): Promise<CBooleanResponse> {
    const { dataSources } = ctx;
    return await dataSources.dataSource.removeInvoice(invoiceId);
  }
}
