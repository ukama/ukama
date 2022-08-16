import Stripe from "stripe";
import { Service } from "typedi";
import { AttachPaymentDto } from "../types";
import { Authentication } from "../../../common/Authentication";
import { Resolver, UseMiddleware, Mutation, Arg } from "type-graphql";
import { STRIP_SK } from "../../../constants";

@Service()
@Resolver()
export class AttachPaymentWithCustomerResolver {
    @Mutation(() => Boolean)
    @UseMiddleware(Authentication)
    async attachPaymentWithCustomer(
        @Arg("data")
        data: AttachPaymentDto
    ): Promise<boolean> {
        const stripe = new Stripe(STRIP_SK, {
            typescript: true,
            apiVersion: "2022-08-01",
        });
        const customer = await stripe.paymentMethods.attach(data.paymentId, {
            customer: data.customerId,
        });

        return customer.id ? true : false;
    }
}
