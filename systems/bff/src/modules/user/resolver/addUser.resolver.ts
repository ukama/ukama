import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { UserService } from "../service";
import { UserInputDto, UserResDto } from "../types";

@Service()
@Resolver()
export class AddUserResolver {
    constructor(private readonly userService: UserService) {}

    @Mutation(() => UserResDto)
    @UseMiddleware(Authentication)
    async addUser(
        @Arg("data") data: UserInputDto,
        @Ctx() ctx: Context
    ): Promise<UserResDto> {
        const user = await this.userService.addUser(data, parseHeaders(ctx));
        return user;
    }
}
