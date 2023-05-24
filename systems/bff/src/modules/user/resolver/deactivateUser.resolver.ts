import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { UserService } from "../service";
import { UserResDto } from "../types";

@Service()
@Resolver()
export class DeactivateUserResolver {
    constructor(private readonly userService: UserService) {}

    @Mutation(() => UserResDto)
    @UseMiddleware(Authentication)
    async deactivateUser(
        @Arg("uuid")
        uuid: string,
        @Ctx() ctx: Context
    ): Promise<UserResDto> {
        return this.userService.deactivateUser(uuid, parseHeaders(ctx));
    }
}
