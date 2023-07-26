import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { SimsResDto } from "../types";
import { SIM_TYPES } from "../../../constants";

@Resolver()
export class GetSimsResolver {

    @Query(() => SimsResDto)
    @UseMiddleware(Authentication)
    async getSims(
        @Arg("type") type: SIM_TYPES,
        @Ctx() ctx: Context
    ): Promise<SimsResDto> {
        const { dataSources } = ctx;
        return dataSources.dataSource.getSims(type, parseHeaders(ctx));
    }
}
