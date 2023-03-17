import { Service } from "typedi";
import { UserService } from "../service";
import { parseCookie } from "../../../common";
import { Context } from "../../../common/types";
import { Authentication } from "../../../common/Authentication";
import { OrgUserSimDto, UpdateUserServiceInput } from "../types";
import { Resolver, Arg, Mutation, UseMiddleware, Ctx } from "type-graphql";

@Service()
@Resolver()
export class UpdateUserRoamingResolver {
    constructor(private readonly userService: UserService) {}

    @Mutation(() => OrgUserSimDto)
    @UseMiddleware(Authentication)
    async updateUserRoaming(
        @Arg("data") data: UpdateUserServiceInput,
        @Ctx() ctx: Context,
    ): Promise<OrgUserSimDto> {
        return this.userService.updateUserRoaming(data, parseCookie(ctx));
    }
}
