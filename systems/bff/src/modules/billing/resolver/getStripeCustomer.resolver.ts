import Stripe from "stripe";
import { Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { STRIP_SK } from "../../../constants";
import { getStripeIdByUserId } from "../../../utils";
import { StripeCustomer } from "../types";

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
            getStripeIdByUserId(parseHeaders(ctx).orgId)
        );
        return {
            id: customer.id,
            name: customer?.name || "name",
            email: customer?.email || "email",
        };
    }
}
