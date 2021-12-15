import { Resolver, Arg, Ctx, Mutation, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { AddUserDto, AddUserResponse } from "../types";
import { UserService } from "../service";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";

@Service()
@Resolver()
export class AddUserResolver {
    constructor(private readonly userService: UserService) {}

    @Mutation(() => AddUserResponse)
    @UseMiddleware(Authentication)
    async addUser(
        @Arg("orgId") orgId: string,
        @Arg("data") data: AddUserDto,
        @Ctx() ctx: Context
    ): Promise<AddUserResponse | null> {
        return await this.userService.addUser(orgId, data, ctx);
    }
}
