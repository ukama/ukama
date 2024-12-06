import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { CorrespondentsInputDto, CorrespondentsResDto } from "./types";

@Resolver()
export class GetCorrespondentsResolver {
  @Query(() => CorrespondentsResDto)
  async getCorrespondents(
    @Arg("data") data: CorrespondentsInputDto,
    @Ctx() ctx: Context
  ): Promise<CorrespondentsResDto> {
    const { dataSources, baseURL } = ctx;
    return dataSources.dataSource.getCorrespondents(
      baseURL,
      data.phoneNumber,
      data.paymentMethod
    );
  }
}
