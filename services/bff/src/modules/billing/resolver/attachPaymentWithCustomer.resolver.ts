import Stripe from "stripe";
import { Service } from "typedi";
import { AttachPaymentDto } from "../types";
import { Authentication } from "../../../common/Authentication";
import { Resolver, UseMiddleware, Mutation, Arg } from "type-graphql";

@Service()
@Resolver()
export class AttachPaymentWithCustomerResolver {
    @Mutation(() => Boolean)
    @UseMiddleware(Authentication)
    async attachPaymentWithCustomer(
        @Arg("data")
        data: AttachPaymentDto
    ): Promise<boolean> {
        const stripe = new Stripe(
            "sk_test_51LN9vGHBOiFTwZOs0zrb6CkarnlRRocpGSoIZa3jL7vtMeolNjrzf7PAL3hMDHQZENnxIvbw8X7Bfx5CxsUfVfyu00HIVQCYAm",
            {
                typescript: true,
                apiVersion: "2022-08-01",
            }
        );
        const customer = await stripe.paymentMethods.attach(data.paymentId, {
            customer: data.customerId,
        });

        return customer.id ? true : false;
    }
}
