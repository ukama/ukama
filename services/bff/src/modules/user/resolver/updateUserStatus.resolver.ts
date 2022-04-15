import { Service } from "typedi";
import { UserService } from "../service";
import { Context } from "../../../common/types";
import { Authentication } from "../../../common/Authentication";
import { UpdateUserServiceInput, UpdateUserServiceRes } from "../types";
import { Resolver, Arg, Mutation, UseMiddleware, Ctx } from "type-graphql";
import { getHeaders } from "../../../common";

@Service()
@Resolver()
export class UpdateUserStatusResolver {
    constructor(private readonly userService: UserService) {}

    @Mutation(() => UpdateUserServiceRes)
    @UseMiddleware(Authentication)
    async updateUserStatus(
        @Arg("data") data: UpdateUserServiceInput,
        @Ctx() ctx: Context
    ): Promise<UpdateUserServiceRes | null> {
        return this.userService.updateUserStatus(data, getHeaders(ctx));
    }
}
