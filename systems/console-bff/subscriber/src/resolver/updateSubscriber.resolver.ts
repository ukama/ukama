import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { BoolResponse, Context } from "../../../common/types";
import { UpdateSubscriberInputDto } from "../types";

@Resolver()
export class UpdateSubscriberResolver {

    @Mutation(() => BoolResponse)
    @UseMiddleware(Authentication)
    async updateSubscriber(
        @Arg("subscriberId") subscriberId: string,
        @Arg("data") data: UpdateSubscriberInputDto,
        @Ctx() ctx: Context
    ): Promise<BoolResponse> {
        const { dataSources } = ctx;
        return await dataSources.dataSource.updateSubscriber(
            subscriberId,
            data,
            parseHeaders(ctx)
        );
    }
}
