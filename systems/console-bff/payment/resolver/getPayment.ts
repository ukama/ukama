import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { PaymentDto } from "./types";

@Resolver()
export class GetPaymentResolver {
  @Query(() => PaymentDto)
  async getPayment(
    @Arg("paymentId") paymentId: string,
    @Ctx() ctx: Context
  ): Promise<PaymentDto> {
    const { dataSources, baseURL } = ctx;
    return dataSources.dataSource.getPayment(baseURL, paymentId);
  }
}
