/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import Stripe from "stripe";
import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { STRIP_SK } from "../../common/configs";
import { getStripeIdByUserId } from "../../common/utils";
import { Context } from "../context";

@Resolver()
export class AttachPaymentWithCustomerResolver {
  @Mutation(() => Boolean)
  async attachPaymentWithCustomer(
    @Arg("paymentId")
    paymentId: string,
    @Ctx() ctx: Context
  ): Promise<boolean> {
    const { headers } = ctx;
    const stripe = new Stripe(STRIP_SK, {
      typescript: true,
      apiVersion: "2024-04-10",
    });
    const customer = await stripe.paymentMethods.attach(paymentId, {
      customer: getStripeIdByUserId(headers.orgId),
    });

    return customer.id ? true : false;
  }
}
