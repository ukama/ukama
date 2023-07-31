import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { GetSimByNetworkInputDto, SimDetailsDto } from "../types";

@Resolver()
export class GetSimByNetworkResolver {

    @Query(() => SimDetailsDto)
    @UseMiddleware(Authentication)
    async getSim(
        @Arg("data") data: GetSimByNetworkInputDto,
        @Ctx() ctx: Context
    ): Promise<SimDetailsDto> {
        const { dataSources } = ctx;
        return await dataSources.dataSource.getSimByNetworkId(data, parseHeaders(ctx));
    }
}
