import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { SubscriberDto } from "../types";

@Resolver()
export class GetSubscriberResolver {

    @Query(() => SubscriberDto)
    @UseMiddleware(Authentication)
    async getSubscriber(
        @Arg("subscriberId") subscriberId: string,
        @Ctx() ctx: Context
    ): Promise<SubscriberDto> {
        const { dataSources } = ctx;
        return await dataSources.dataSource.getSubscriber(
            subscriberId,
            parseHeaders(ctx)
        );
    }
}
