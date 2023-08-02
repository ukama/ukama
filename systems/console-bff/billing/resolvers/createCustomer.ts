import Stripe from "stripe";
import { Arg, Mutation, Resolver } from "type-graphql";

import { STRIP_SK } from "../../common/configs";
import { CreateCustomerDto, StripeCustomer } from "./types";

@Resolver()
export class CreateCustomerResolver {
  @Mutation(() => StripeCustomer)
  async createCustomer(
    @Arg("data")
    req: CreateCustomerDto
  ): Promise<StripeCustomer> {
    const stripe = new Stripe(STRIP_SK, {
      typescript: true,
      apiVersion: "2022-11-15",
    });
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
