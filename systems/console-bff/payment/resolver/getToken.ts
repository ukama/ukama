import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { TokenResDto } from "./types";

@Resolver()
export class GetTokenResolver {
  @Query(() => TokenResDto)
  async getToken(
    @Arg("paymentId") paymentId: string,
    @Ctx() ctx: Context
  ): Promise<TokenResDto> {
    const { dataSources, baseURL } = ctx;
    return dataSources.dataSource.getToken(baseURL, paymentId);
  }
}
