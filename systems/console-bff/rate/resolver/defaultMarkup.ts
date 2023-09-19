import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { CBooleanResponse } from "../../common/types";
import { Context } from "../context";
import { DefaultMarkupInputDto } from "./types";

@Resolver()
export class DefaultMarkupResolver {
  @Mutation(() => CBooleanResponse)
  @UseMiddleware(Authentication)
  async defaultMarkup(
    @Arg("data") data: DefaultMarkupInputDto,
    @Ctx() ctx: Context
  ): Promise<CBooleanResponse> {
    const { dataSources } = ctx;
    return dataSources.dataSource.defaultMarkup(data);
  }
}
