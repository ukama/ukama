import Stripe from "stripe";
import { Service } from "typedi";
import { CreateCustomerDto, StripeCustomer } from "../types";
import { Authentication } from "../../../common/Authentication";
import { Resolver, UseMiddleware, Mutation, Arg } from "type-graphql";

@Service()
@Resolver()
export class CreateCustomerResolver {
    @Mutation(() => StripeCustomer)
    @UseMiddleware(Authentication)
    async createCustomer(
        @Arg("data")
        req: CreateCustomerDto
    ): Promise<StripeCustomer> {
        const stripe = new Stripe(
            "sk_test_51LN9vGHBOiFTwZOs0zrb6CkarnlRRocpGSoIZa3jL7vtMeolNjrzf7PAL3hMDHQZENnxIvbw8X7Bfx5CxsUfVfyu00HIVQCYAm",
            {
                typescript: true,
                apiVersion: "2022-08-01",
            }
        );
        const customer = await stripe.customers.create({
            name: req.name,
            email: req.email,
            description: "Test ukama customer",
        });

        return {
            id: customer.id,
            name: customer?.name || "name",
            email: customer?.email || "email",
        };
    }
}
