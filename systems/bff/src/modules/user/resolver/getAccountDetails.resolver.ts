import { Resolver, Query, UseMiddleware, Ctx } from "type-graphql";
import { Service } from "typedi";
import { UserService } from "../service";
import { GetAccountDetailsDto } from "../types";
import { Authentication } from "../../../common/Authentication";
import { parseCookie } from "../../../common";
import { Context } from "../../../common/types";

@Service()
@Resolver()
export class GetAccountDetailsResolver {
    constructor(private readonly userService: UserService) {}

    @Query(() => GetAccountDetailsDto)
    @UseMiddleware(Authentication)
    async getAccountDetails(
        @Ctx() ctx: Context,
    ): Promise<GetAccountDetailsDto | null> {
        return this.userService.getAccountDetails(parseCookie(ctx));
    }
}
