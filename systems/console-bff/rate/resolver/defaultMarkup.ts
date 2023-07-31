import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { BoolResponse } from "../../common/types";
import { Context } from "../context";
import { DefaultMarkupInputDto } from "./types";

@Resolver()
export class DefaultMarkupResolver {
  @Mutation(() => BoolResponse)
  @UseMiddleware(Authentication)
  async defaultMarkup(
    @Arg("data") data: DefaultMarkupInputDto,
    @Ctx() ctx: Context
  ): Promise<BoolResponse> {
    const { dataSources } = ctx;
    return dataSources.dataSource.defaultMarkup(data);
  }
}
