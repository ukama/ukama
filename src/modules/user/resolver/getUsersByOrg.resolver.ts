import { Resolver, Query, UseMiddleware, Ctx, Arg } from "type-graphql";
import { Service } from "typedi";
import { OrgUserResponseDto } from "../types";
import { UserService } from "../service";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";

@Service()
@Resolver()
export class GetUsersByOrgResolver {
    constructor(private readonly userService: UserService) {}

    @Query(() => OrgUserResponseDto)
    @UseMiddleware(Authentication)
    async myUsers(
        @Arg("orgId") orgId: string,
        @Ctx() ctx: Context
    ): Promise<OrgUserResponseDto> {
        return this.userService.getUsersByOrg(orgId, ctx);
    }
}
