import { Resolver, Arg, Mutation, UseMiddleware, Ctx } from "type-graphql";
import { Service } from "typedi";
import { DeactivateResponse } from "../types";
import { UserService } from "../service";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { parseCookie } from "../../../common";

@Service()
@Resolver()
export class DeactivateUserResolver {
    constructor(private readonly userService: UserService) {}

    @Mutation(() => DeactivateResponse)
    @UseMiddleware(Authentication)
    async deactivateUser(
        @Arg("id")
        id: string,
        @Ctx() ctx: Context
    ): Promise<DeactivateResponse> {
        return this.userService.deactivateUser(id, parseCookie(ctx));
    }
}
