import { Resolver, Arg, Mutation, UseMiddleware, Ctx } from "type-graphql";
import { Service } from "typedi";
import { UserInputDto, UserResDto } from "../types";
import { UserService } from "../service";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { parseCookie } from "../../../common";

@Service()
@Resolver()
export class UpdateUserResolver {
    constructor(private readonly userService: UserService) {}

    @Mutation(() => UserResDto)
    @UseMiddleware(Authentication)
    async updateUser(
        @Arg("userId") userId: string,
        @Arg("data") data: UserInputDto,
        @Ctx() ctx: Context
    ): Promise<UserResDto | null> {
        return this.userService.updateUser(userId, data, parseCookie(ctx));
    }
}
