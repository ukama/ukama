import { Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { BillResponse } from "./types";

@Resolver()
export class GetCurrentBillResolver {
  @Query(() => BillResponse)
  async getCurrentBill(@Ctx() ctx: Context): Promise<BillResponse> {
    const { dataSources } = ctx;
    return dataSources.dataSource.getCurrentBill();
  }
}
