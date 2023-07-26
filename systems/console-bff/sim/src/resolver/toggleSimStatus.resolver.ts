import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { SimStatusResDto, ToggleSimStatusInputDto } from "../types";

@Resolver()
export class ToggleSimStatusResolver {

    @Mutation(() => SimStatusResDto)
    @UseMiddleware(Authentication)
    async toggleSimStatus(
        @Arg("data") data: ToggleSimStatusInputDto,
        @Ctx() ctx: Context
    ): Promise<SimStatusResDto> {
        const { dataSources } = ctx;
        return await dataSources.dataSource.toggleSimStatus(data, parseHeaders(ctx));
    }
}
