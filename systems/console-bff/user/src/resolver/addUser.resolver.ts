import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";

import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { UserInputDto, UserResDto } from "../types";

@Resolver()
export class AddUserResolver {
    constructor(private readonly userService: UserService) {}

    @Mutation(() => UserResDto)
    @UseMiddleware(Authentication)
    async addUser(
        @Arg("data") data: UserInputDto,
        @Ctx() ctx: Context
    ): Promise<UserResDto> {
        const { dataSources } = ctx;
        const user = await dataSources.dataSource.addUser(data, parseHeaders(ctx));
        return user;
    }
}
