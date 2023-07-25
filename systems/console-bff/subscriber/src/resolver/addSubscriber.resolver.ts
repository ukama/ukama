import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { SubscriberDto, SubscriberInputDto } from "../types";

@Resolver()
export class AddSubscriberResolver {
    @Mutation(() => SubscriberDto)
    @UseMiddleware(Authentication)
    async addSubscriber(
        @Arg("data") data: SubscriberInputDto,
        @Ctx() ctx: Context
    ): Promise<SubscriberDto> {
        const { dataSources } = ctx;
        return await dataSources.dataSource.addSubscriber(data, parseHeaders(ctx));
    }
}
