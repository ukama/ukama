import { Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { BillHistoryDto } from "./types";

@Resolver()
export class GetBillHistoryResolver {
  @Query(() => [BillHistoryDto])
  async getBillHistory(@Ctx() ctx: Context): Promise<BillHistoryDto[]> {
    const { dataSources } = ctx;
    return dataSources.dataSource.getBillHistory();
  }
}
