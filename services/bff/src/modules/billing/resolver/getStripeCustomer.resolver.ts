import Stripe from "stripe";
import { Service } from "typedi";
import { StripeCustomer } from "../types";
import { Resolver, Query, UseMiddleware, Arg } from "type-graphql";
import { Authentication } from "../../../common/Authentication";
import { STRIP_SK } from "../../../constants";

@Service()
@Resolver()
export class GetStripeCustomerResolver {
    @Query(() => StripeCustomer)
    @UseMiddleware(Authentication)
    async getStripeCustomer(
        @Arg("id")
        id: string
    ): Promise<StripeCustomer> {
        const stripe = new Stripe(STRIP_SK, {
            typescript: true,
            apiVersion: "2022-08-01",
        });
        const customer: any = await stripe.customers.retrieve(id);
        return {
            id: customer.id,
            name: customer?.name || "name",
            email: customer?.email || "email",
        };
    }
}
