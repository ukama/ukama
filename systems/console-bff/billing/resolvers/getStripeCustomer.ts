import Stripe from "stripe";
import { Ctx, Query, Resolver } from "type-graphql";

import { STRIP_SK } from "../../common/configs";
import { getStripeIdByUserId, parseHeaders } from "../../common/utils";
import { Context } from "../context";
import { StripeCustomer } from "./types";

@Resolver()
export class GetStripeCustomerResolver {
  @Query(() => StripeCustomer)
  async getStripeCustomer(@Ctx() ctx: Context): Promise<StripeCustomer> {
    const stripe = new Stripe(STRIP_SK, {
      typescript: true,
      apiVersion: "2022-11-15",
    });
    const customer: any = await stripe.customers.retrieve(
      getStripeIdByUserId(parseHeaders().orgId)
    );
    return {
      id: customer.id,
      name: customer?.name || "name",
      email: customer?.email || "email",
    };
  }
}
