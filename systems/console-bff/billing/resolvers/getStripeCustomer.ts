/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import Stripe from "stripe";
import { Ctx, Query, Resolver } from "type-graphql";

import { STRIP_SK } from "../../common/configs";
import { getStripeIdByUserId } from "../../common/utils";
import { Context } from "../context";
import { StripeCustomer } from "./types";

@Resolver()
export class GetStripeCustomerResolver {
  @Query(() => StripeCustomer)
  async getStripeCustomer(@Ctx() ctx: Context): Promise<StripeCustomer> {
    const { headers } = ctx;
    const stripe = new Stripe(STRIP_SK, {
      typescript: true,
      apiVersion: "2022-11-15",
    });
    const customer: any = await stripe.customers.retrieve(
      getStripeIdByUserId(headers.orgId)
    );
    return {
      id: customer.id,
      name: customer?.name || "name",
      email: customer?.email || "email",
    };
  }
}
