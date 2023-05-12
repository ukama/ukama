import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseCookie } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { MemberObj } from "../types";
import { OrgService } from "./../service";

@Service()
@Resolver()
export class AddMemberResolver {
    constructor(private readonly OrgService: OrgService) {}

    @Mutation(() => MemberObj)
    @UseMiddleware(Authentication)
    async addMember(
        @Arg("userId") userId: string,
        @Ctx() ctx: Context
    ): Promise<MemberObj> {
        return this.OrgService.addMember(userId, parseCookie(ctx));
    }
}
