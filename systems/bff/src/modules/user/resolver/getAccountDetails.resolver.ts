import { Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { UserService } from "../service";
import { GetAccountDetailsDto } from "../types";

@Service()
@Resolver()
export class GetAccountDetailsResolver {
    constructor(private readonly userService: UserService) {}

    @Query(() => GetAccountDetailsDto)
    @UseMiddleware(Authentication)
    async getAccountDetails(
        @Ctx() ctx: Context
    ): Promise<GetAccountDetailsDto | null> {
        return this.userService.getAccountDetails(parseHeaders(ctx));
    }
}
