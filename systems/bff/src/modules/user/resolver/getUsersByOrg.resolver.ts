import { Service } from "typedi";
import { GetUsersDto } from "../types";
import { UserService } from "../service";
import { parseCookie } from "../../../common";
import { Context } from "../../../common/types";
import { Authentication } from "../../../common/Authentication";
import { Resolver, Query, UseMiddleware, Ctx } from "type-graphql";

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
