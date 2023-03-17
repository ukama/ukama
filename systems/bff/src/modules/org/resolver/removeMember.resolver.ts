import { Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseCookie } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { BoolResponse, Context } from "../../../common/types";
import { OrgService } from "./../service";

@Service()
@Resolver()
export class RemoveMemberResolver {
    constructor(private readonly OrgService: OrgService) {}

    @Mutation(() => BoolResponse)
    @UseMiddleware(Authentication)
    async removeMember(@Ctx() ctx: Context): Promise<BoolResponse> {
        return this.OrgService.removeMember(parseCookie(ctx));
    }
}
