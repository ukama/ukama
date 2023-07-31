import { Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { DefaultMarkupHistoryResDto } from "../types";

@Resolver()
export class GetDefaultMarkupHistoryResolver {
    @Query(() => DefaultMarkupHistoryResDto)
    @UseMiddleware(Authentication)
    async getDefaultMarkupHistory(
        @Ctx() ctx: Context
    ): Promise<DefaultMarkupHistoryResDto> {
        const { dataSources } = ctx;
        return dataSources.dataSource.getDefaultMarkupHistory(parseHeaders(ctx));
    }
}
