import { Resolver, Arg, Ctx, Mutation, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { ActivateUserResponse } from "../types";
import { UserService } from "../service";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { getHeaders } from "../../../common";

@Service()
@Resolver()
export class AddUserResolver {
    constructor(private readonly userService: UserService) {}

    @Mutation(() => ActivateUserResponse)
    @UseMiddleware(Authentication)
    async deleteUser(
        @Arg("orgId") orgId: string,
        @Arg("userId") userId: string,
        @Ctx() ctx: Context
    ): Promise<ActivateUserResponse | null> {
        return this.userService.deleteUser(orgId, userId, getHeaders(ctx));
    }
}
