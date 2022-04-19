import { Resolver, Arg, Ctx, Mutation, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { AddUserDto, UserResDto } from "../types";
import { UserService } from "../service";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { parseCookie } from "../../../common";

@Service()
@Resolver()
export class AddUserResolver {
    constructor(private readonly userService: UserService) {}

    @Mutation(() => UserResDto)
    @UseMiddleware(Authentication)
    async addUser(
        @Arg("data") data: AddUserDto,
        @Ctx() ctx: Context
    ): Promise<UserResDto | null> {
        return this.userService.addUser(data, parseCookie(ctx));
    }
}
