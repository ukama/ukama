import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { BoolResponse, Context } from "../../../common/types";

@Resolver()
export class DeleteSubscriberResolver {

    @Mutation(() => BoolResponse)
    @UseMiddleware(Authentication)
    async deleteSubscriber(
        @Arg("subscriberId") subscriberId: string,
        @Ctx() ctx: Context
    ): Promise<BoolResponse> {
        const { dataSources } = ctx;
        return await dataSources.dataSource.deleteSubscriber(
            subscriberId,
            parseHeaders(ctx)
        );
    }
}
