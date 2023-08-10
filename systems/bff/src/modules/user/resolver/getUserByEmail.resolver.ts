import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { UserService } from "../service";
import { UserResDto } from "../types";

@Service()
@Resolver()
export class GetUserByEmailResolver {
    constructor(private readonly userService: UserService) {}

    @Query(() => UserResDto)
    @UseMiddleware(Authentication)
    async getUserByEmail(
        @Arg("email") email: string,
        @Ctx() ctx: Context
    ): Promise<UserResDto | null> {
        return this.userService.getUserByEmail(email, parseHeaders(ctx));
    }
}
