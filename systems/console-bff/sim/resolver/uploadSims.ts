import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { Context } from "../context";
import { UploadSimsInputDto, UploadSimsResDto } from "./types";

@Resolver()
export class UploadSimsResolver {
  @Mutation(() => UploadSimsResDto)
  @UseMiddleware(Authentication)
  async uploadSims(
    @Arg("data") data: UploadSimsInputDto,
    @Ctx() ctx: Context
  ): Promise<UploadSimsResDto> {
    const { dataSources } = ctx;
    return await dataSources.dataSource.uploadSims(data);
  }
}
