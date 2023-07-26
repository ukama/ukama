import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { AllocateSimInputDto, SimDto } from "../types";

@Resolver()
export class AllocateSimResolver {

    @Mutation(() => SimDto)
    @UseMiddleware(Authentication)
    async allocateSim(
        @Arg("data") data: AllocateSimInputDto,
        @Ctx() ctx: Context
    ): Promise<SimDto> {
        const { dataSources } = ctx;
        return await dataSources.dataSource.allocateSim(data, parseHeaders(ctx));
    }
}
