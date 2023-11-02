import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { SIM_TYPES } from "../../common/enums";
import { Context } from "../context";
import { SimsResDto } from "./types";

@Resolver()
export class GetSimsResolver {
  @Query(() => SimsResDto)
  async getSims(
    @Arg("type") type: SIM_TYPES,
    @Ctx() ctx: Context
  ): Promise<SimsResDto> {
    const { dataSources } = ctx;
    return dataSources.dataSource.getSims(type);
  }
}
