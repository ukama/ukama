import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { AddOrgInputDto, OrgDto } from "../types";
import { OrgService } from "./../service";

@Service()
@Resolver()
export class AddOrgResolver {
    constructor(private readonly OrgService: OrgService) {}

    @Mutation(() => OrgDto)
    @UseMiddleware(Authentication)
    async addOrg(
        @Arg("data") data: AddOrgInputDto,
        @Ctx() ctx: Context
    ): Promise<OrgDto> {
        return this.OrgService.addOrg(data, parseHeaders(ctx));
    }
}
