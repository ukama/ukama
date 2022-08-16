import Stripe from "stripe";
import { Service } from "typedi";
import { STRIP_SK } from "../../../constants";
import { StripePaymentMethods } from "../types";
import { Authentication } from "../../../common/Authentication";
import { Resolver, Query, UseMiddleware, Arg } from "type-graphql";

@Service()
@Resolver()
export class RetrivePaymentMethodsResolver {
    @Query(() => [StripePaymentMethods])
    @UseMiddleware(Authentication)
    async retrivePaymentMethods(
        @Arg("id")
        id: string
    ): Promise<StripePaymentMethods[]> {
        const stripe = new Stripe(STRIP_SK, {
            typescript: true,
            apiVersion: "2022-08-01",
        });
        const pm: Stripe.ApiList<Stripe.PaymentMethod> =
            await stripe.customers.listPaymentMethods(id, {
                type: "card",
            });
        const list: StripePaymentMethods[] = [];
        for (const ele of pm.data) {
            if (ele.card) {
                list.push({
                    id: ele.id,
                    type: ele.type,
                    created: ele.created,
                    brand: ele.card?.brand,
                    last4: ele.card?.last4,
                    funding: ele.card?.funding,
                    exp_year: ele.card?.exp_year,
                    exp_month: ele.card?.exp_month,
                    country: ele.card?.country || undefined,
                    cvc_check: ele.card?.checks?.cvc_check || undefined,
                });
            } else {
                list.push({
                    brand: "",
                    last4: "",
                    id: ele.id,
                    funding: "",
                    exp_year: 0,
                    exp_month: 0,
                    type: ele.type,
                    country: undefined,
                    created: ele.created,
                    cvc_check: undefined,
                });
            }
        }
        return list;
    }
}
