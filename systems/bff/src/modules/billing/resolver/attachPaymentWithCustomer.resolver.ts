import Stripe from "stripe";
import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { STRIP_SK } from "../../../constants";
import { getStripeIdByUserId } from "../../../utils";

@Service()
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
            apiVersion: "2022-08-01",
        });
        const customer = await stripe.paymentMethods.attach(paymentId, {
            customer: getStripeIdByUserId(parseHeaders(ctx).orgId),
        });

        return customer.id ? true : false;
    }
}
