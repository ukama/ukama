import { Ctx, Query, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/Authentication";
import { Context } from "../context";
import { BillHistoryDto } from "./types";

@Resolver()
export class GetBillHistoryResolver {
  @Query(() => [BillHistoryDto])
  @UseMiddleware(Authentication)
  async getBillHistory(@Ctx() ctx: Context): Promise<BillHistoryDto[]> {
    const { dataSources } = ctx;
    return dataSources.dataSource.getBillHistory();
  }
}
