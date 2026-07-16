/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { logger } from "../../common/logger";
import type { AppContext } from "../../server/context";
import { AddPaymentInputDto, PaymentDto } from "./types";

@Resolver()
export class AddPaymentResolver {
  @Mutation(() => PaymentDto)
  async addPayment(
    @Arg("data") data: AddPaymentInputDto,
    @Ctx() ctx: AppContext
  ): Promise<PaymentDto> {
    const { dataSources } = ctx;
    const baseURL = await ctx.urls.url("payments");
    if (
      data.paymentMethod?.toLowerCase() === "cash" &&
      data.itemType?.toLowerCase() === "package" &&
      !data.sim
    ) {
      logger.error("Missing sim for cash package payment");
      throw new Error("sim is required for a cash package payment");
    }

    logger.info(
      `Adding ${data.paymentMethod} payment for ${data.itemType}: ${data.itemId}`
    );

    return dataSources.payment.addPayment(baseURL, data);
  }
}
