import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { UserService } from "../service";
import { UserResDto } from "../types";

@Service()
@Resolver()
export class GetUserResolver {
    constructor(private readonly userService: UserService) {}

    @Query(() => UserResDto)
    @UseMiddleware(Authentication)
    async getUser(
        @Arg("userId") userId: string,
        @Ctx() ctx: Context
    ): Promise<UserResDto | null> {
        return this.userService.getUser(userId, parseHeaders(ctx));
    }
}
