import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { AddMemberInputDto, MemberObj } from "../types";
import { OrgService } from "./../service";

@Service()
@Resolver()
export class AddMemberResolver {
    constructor(private readonly OrgService: OrgService) {}

    @Mutation(() => MemberObj)
    @UseMiddleware(Authentication)
    async addMember(
        @Arg("data") data: AddMemberInputDto,
        @Ctx() ctx: Context
    ): Promise<MemberObj> {
        return this.OrgService.addMember(data, parseHeaders(ctx));
    }
}
