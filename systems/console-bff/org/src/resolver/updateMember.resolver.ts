import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { BoolResponse, Context } from "../../../common/types";
import { UpdateMemberInputDto } from "../types";

@Resolver()
export class UpdateMemberResolver {

    @Mutation(() => BoolResponse)
    @UseMiddleware(Authentication)
    async updateMember(
        @Arg("memberId") memberId: string,
        @Arg("data") data: UpdateMemberInputDto,
        @Ctx() ctx: Context
    ): Promise<BoolResponse> {
        const { dataSources } = ctx;
        return dataSources.dataSource.updateMember(memberId, data, parseHeaders(ctx));
    }
}
