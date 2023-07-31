import { Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { DefaultMarkupResDto } from "../types";

@Resolver()
export class GetDefaultMarkupResolver {

    @Query(() => DefaultMarkupResDto)
    @UseMiddleware(Authentication)
    async getDefaultMarkup(@Ctx() ctx: Context): Promise<DefaultMarkupResDto> {
        const { dataSources } = ctx;
        return dataSources.dataSource.getDefaultMarkup(parseHeaders(ctx));
    }
}
