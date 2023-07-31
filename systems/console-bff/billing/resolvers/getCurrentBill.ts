import { Ctx, Query, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/Authentication";
import { Context } from "../context";
import { BillResponse } from "./types";

@Resolver()
export class GetCurrentBillResolver {
  @Query(() => BillResponse)
  @UseMiddleware(Authentication)
  async getCurrentBill(@Ctx() ctx: Context): Promise<BillResponse> {
    const { dataSources } = ctx;
    return dataSources.dataSource.getCurrentBill();
  }
}
