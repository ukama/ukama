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
      apiVersion: "2022-11-15",
    });
    const customer = await stripe.paymentMethods.attach(paymentId, {
      customer: getStripeIdByUserId(headers.orgId),
    });

    return customer.id ? true : false;
  }
}
