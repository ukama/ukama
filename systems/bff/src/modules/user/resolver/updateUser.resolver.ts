import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { UserService } from "../service";
import { UpdateUserInputDto, UserResDto } from "../types";

@Service()
@Resolver()
export class UpdateUserResolver {
    constructor(private readonly userService: UserService) {}

    @Mutation(() => UserResDto)
    @UseMiddleware(Authentication)
    async updateUser(
        @Arg("userId") userId: string,
        @Arg("data") data: UpdateUserInputDto,
        @Ctx() ctx: Context
    ): Promise<UserResDto | null> {
        return this.userService.updateUser(userId, data, parseHeaders(ctx));
    }
}
