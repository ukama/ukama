import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { BoolResponse, Context } from "../../../common/types";
import { UpdateMemberInputDto } from "../types";
import { OrgService } from "./../service";

@Service()
@Resolver()
export class updateMemberResolver {
    constructor(private readonly OrgService: OrgService) {}

    @Mutation(() => BoolResponse)
    @UseMiddleware(Authentication)
    async updateMember(
        @Arg("memberId") memberId: string,
        @Arg("data") data: UpdateMemberInputDto,
        @Ctx() ctx: Context
    ): Promise<BoolResponse> {
        return this.OrgService.updateMember(memberId, data, parseHeaders(ctx));
    }
}
