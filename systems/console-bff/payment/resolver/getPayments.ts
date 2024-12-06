import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { GetPaymentsInputDto, PaymentsDto } from "./types";

@Resolver()
export class GetPaymentsResolver {
  @Query(() => PaymentsDto)
  async getPayments(
    @Arg("data") data: GetPaymentsInputDto,
    @Ctx() ctx: Context
  ): Promise<PaymentsDto> {
    const { dataSources, baseURL } = ctx;
    return dataSources.dataSource.getPayments(baseURL, data);
  }
}
