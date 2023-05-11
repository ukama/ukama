import Stripe from "stripe";
import { Service } from "typedi";
import { StripeCustomer } from "../types";
import { STRIP_SK } from "../../../constants";
import { parseCookie } from "../../../common";
import { Context } from "../../../common/types";
import { getStripeIdByUserId } from "../../../utils";
import { Authentication } from "../../../common/Authentication";
import { Resolver, Query, UseMiddleware, Ctx } from "type-graphql";

@Service()
@Resolver()
export class GetStripeCustomerResolver {
    @Query(() => StripeCustomer)
    @UseMiddleware(Authentication)
    async getStripeCustomer(@Ctx() ctx: Context): Promise<StripeCustomer> {
        const stripe = new Stripe(STRIP_SK, {
            typescript: true,
            apiVersion: "2022-08-01",
        });
        const customer: any = await stripe.customers.retrieve(
            getStripeIdByUserId(parseCookie(ctx).orgId),
        );
        return {
            id: customer.id,
            name: customer?.name || "name",
            email: customer?.email || "email",
        };
    }
}
