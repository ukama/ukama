import { Service } from "typedi";
import { UserService } from "../service";
import { Context } from "../../../common/types";
import { Authentication } from "../../../common/Authentication";
import { UpdateUserServiceInput, OrgUserSimDto } from "../types";
import { Resolver, Arg, Mutation, UseMiddleware, Ctx } from "type-graphql";
import { parseCookie } from "../../../common";

@Service()
@Resolver()
export class UpdateUserStatusResolver {
    constructor(private readonly userService: UserService) {}

    @Mutation(() => OrgUserSimDto)
    @UseMiddleware(Authentication)
    async updateUserStatus(
        @Arg("data") data: UpdateUserServiceInput,
        @Ctx() ctx: Context,
    ): Promise<OrgUserSimDto> {
        return this.userService.updateUserStatus(data, parseCookie(ctx));
    }
}
