import Stripe from "stripe";
import { Service } from "typedi";
import { STRIP_SK } from "../../../constants";
import { parseCookie } from "../../../common";
import { Context } from "../../../common/types";
import { getStripeIdByUserId } from "../../../utils";
import { Authentication } from "../../../common/Authentication";
import { Resolver, UseMiddleware, Mutation, Arg, Ctx } from "type-graphql";

@Service()
@Resolver()
export class AttachPaymentWithCustomerResolver {
    @Mutation(() => Boolean)
    @UseMiddleware(Authentication)
    async attachPaymentWithCustomer(
        @Arg("paymentId")
        paymentId: string,
        @Ctx() ctx: Context,
    ): Promise<boolean> {
        const stripe = new Stripe(STRIP_SK, {
            typescript: true,
            apiVersion: "2022-08-01",
        });
        const customer = await stripe.paymentMethods.attach(paymentId, {
            customer: getStripeIdByUserId(parseCookie(ctx).orgId),
        });

        return customer.id ? true : false;
    }
}
