import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { AddInvitationInputDto, InvitationDto } from "../types";
import { OrgService } from "./../service";

@Service()
@Resolver()
export class AddInvitationResolver {
    constructor(private readonly OrgService: OrgService) {}

    @Mutation(() => InvitationDto)
    @UseMiddleware(Authentication)
    async addInvitation(
        @Arg("data") data: AddInvitationInputDto,
        @Ctx() ctx: Context
    ): Promise<InvitationDto> {
        return this.OrgService.addInvitation(data, parseHeaders(ctx));
    }
}
