import { Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { BoolResponse, Context } from "../../../common/types";

@Resolver()
export class RemoveMemberResolver {

    @Mutation(() => BoolResponse)
    @UseMiddleware(Authentication)
    async removeMember(@Ctx() ctx: Context): Promise<BoolResponse> {
        const { dataSources } = ctx;
        return dataSources.dataSource.removeMember(parseHeaders(ctx));
    }
}
