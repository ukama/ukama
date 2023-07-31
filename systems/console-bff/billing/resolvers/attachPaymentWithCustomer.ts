import Stripe from "stripe";
import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";

import { parseHeaders } from "../../common";
import { Authentication } from "../../common/Authentication";
import { STRIP_SK } from "../../common/configs";
import { Context } from "../../common/types";
import { getStripeIdByUserId } from "../../utils";

@Resolver()
export class AttachPaymentWithCustomerResolver {
  @Mutation(() => Boolean)
  @UseMiddleware(Authentication)
  async attachPaymentWithCustomer(
    @Arg("paymentId")
    paymentId: string,
    @Ctx() ctx: Context
  ): Promise<boolean> {
    const stripe = new Stripe(STRIP_SK, {
      typescript: true,
      apiVersion: "2022-11-15",
    });
    const customer = await stripe.paymentMethods.attach(paymentId, {
      customer: getStripeIdByUserId(parseHeaders(ctx).orgId),
    });

    return customer.id ? true : false;
  }
}
