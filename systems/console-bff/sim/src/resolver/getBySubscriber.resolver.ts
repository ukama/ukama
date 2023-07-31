import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { GetSimBySubscriberIdInputDto, SimDetailsDto } from "../types";

@Resolver()
export class GetSimBySubscriberResolver {

    @Query(() => SimDetailsDto)
    @UseMiddleware(Authentication)
    async getSim(
        @Arg("data") data: GetSimBySubscriberIdInputDto,
        @Ctx() ctx: Context
    ): Promise<SimDetailsDto> {
        const { dataSources } = ctx;
        return await dataSources.dataSource.getSimBySubscriberId(
            data,
            parseHeaders(ctx)
        );
    }
}
