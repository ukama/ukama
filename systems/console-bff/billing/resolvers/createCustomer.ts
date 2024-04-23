/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import Stripe from "stripe";
import { Arg, Mutation, Resolver } from "type-graphql";

import { STRIP_SK } from "../../common/configs";
import { CreateCustomerDto, StripeCustomer } from "./types";

@Resolver()
export class CreateCustomerResolver {
  @Mutation(() => StripeCustomer)
  async createCustomer(
    @Arg("data")
    req: CreateCustomerDto
  ): Promise<StripeCustomer> {
    const stripe = new Stripe(STRIP_SK, {
      typescript: true,
      apiVersion: "2024-04-10",
    });
    const customer = await stripe.customers.create({
      name: req.name,
      email: req.email,
      description: "Test ukama customer",
    });

    return {
      id: customer.id,
      name: customer?.name ?? "name",
      email: customer?.email ?? "email",
    };
  }
}
