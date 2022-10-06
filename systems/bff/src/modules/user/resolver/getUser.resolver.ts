import { Resolver, Arg, Query, UseMiddleware, Ctx } from "type-graphql";
import { Service } from "typedi";
import { UserService } from "../service";
import { GetUserDto } from "../types";
import { Authentication } from "../../../common/Authentication";
import { parseCookie } from "../../../common";
import { Context } from "../../../common/types";

@Service()
@Resolver()
export class GetUserResolver {
    constructor(private readonly userService: UserService) {}

    @Query(() => GetUserDto)
    @UseMiddleware(Authentication)
    async getUser(
        @Arg("userId") userId: string,
        @Ctx() ctx: Context
    ): Promise<GetUserDto | null> {
        return this.userService.getUser(userId, parseCookie(ctx));
    }
}
