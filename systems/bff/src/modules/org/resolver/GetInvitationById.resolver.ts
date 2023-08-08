import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { OrgService } from "../service";
import { InvitationResDto } from "../types";

@Service()
@Resolver()
export class GetInvitationByIdResolver {
    constructor(private readonly orgService: OrgService) {}
    @Query(() => InvitationResDto) // Update the return type to InvitationResDto[]
    @UseMiddleware(Authentication)
    async getInvitationById(
        @Arg("id") id: string,
        @Ctx() ctx: Context
    ): Promise<InvitationResDto> {
        return this.orgService.getInvitationById(id, parseHeaders(ctx));
    }
}
