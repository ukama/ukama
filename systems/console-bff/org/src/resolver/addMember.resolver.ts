import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { AddMemberInputDto, MemberObj } from "../types";

@Resolver()
export class AddMemberResolver {

    @Mutation(() => MemberObj)
    @UseMiddleware(Authentication)
    async addMember(
        @Arg("data") data: AddMemberInputDto,
        @Ctx() ctx: Context
    ): Promise<MemberObj> {
        const { dataSources } = ctx;
        return dataSources.dataSource.addMember(data, parseHeaders(ctx));
    }
}
