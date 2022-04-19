import { Resolver, Query, UseMiddleware, Ctx, Arg } from "type-graphql";
import { Service } from "typedi";
import { GetUsersDto } from "../types";
import { UserService } from "../service";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { getHeaders } from "../../../common";

@Service()
@Resolver()
export class GetUsersByOrgResolver {
    constructor(private readonly userService: UserService) {}

    @Query(() => [GetUsersDto])
    @UseMiddleware(Authentication)
    async getUsersByOrg(
        @Arg("orgId") orgId: string,
        @Ctx() ctx: Context
    ): Promise<GetUsersDto[]> {
        return this.userService.getUsersByOrg(orgId, getHeaders(ctx));
    }
}
