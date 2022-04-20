import { Resolver, Query, UseMiddleware, Ctx } from "type-graphql";
import { Service } from "typedi";
import { GetUsersDto } from "../types";
import { UserService } from "../service";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { parseCookie } from "../../../common";

@Service()
@Resolver()
export class GetUsersByOrgResolver {
    constructor(private readonly userService: UserService) {}

    @Query(() => [GetUsersDto])
    @UseMiddleware(Authentication)
    async getUsersByOrg(@Ctx() ctx: Context): Promise<GetUsersDto[]> {
        return this.userService.getUsersByOrg(parseCookie(ctx));
    }
}
