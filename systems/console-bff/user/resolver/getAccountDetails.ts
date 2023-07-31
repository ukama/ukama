import { Ctx, Query, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { parseHeaders } from "../../common/utils";
import { Context } from "../context";
import { GetAccountDetailsDto } from "./types";

@Resolver()
export class GetAccountDetailsResolver {
  @Query(() => GetAccountDetailsDto)
  @UseMiddleware(Authentication)
  async getAccountDetails(
    @Ctx() ctx: Context
  ): Promise<GetAccountDetailsDto | null> {
    const { dataSources } = ctx;
    return dataSources.dataSource.getAccountDetails(parseHeaders());
  }
}
