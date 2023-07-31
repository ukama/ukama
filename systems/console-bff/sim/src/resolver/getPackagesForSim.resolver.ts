import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { GetPackagesForSimInputDto, GetPackagesForSimResDto } from "../types";

@Resolver()
export class GetPackagesForSimResolver {

    @Query(() => GetPackagesForSimResDto)
    @UseMiddleware(Authentication)
    async getSim(
        @Arg("data") data: GetPackagesForSimInputDto,
        @Ctx() ctx: Context
    ): Promise<GetPackagesForSimResDto> {
        const { dataSources } = ctx;

        return await dataSources.dataSource.getPackagesForSim(data, parseHeaders(ctx));
    }
}
