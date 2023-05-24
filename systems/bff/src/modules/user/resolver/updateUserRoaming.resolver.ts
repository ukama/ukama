import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { UserService } from "../service";
import { OrgUserSimDto, UpdateUserServiceInput } from "../types";

@Service()
@Resolver()
export class UpdateUserRoamingResolver {
    constructor(private readonly userService: UserService) {}

    @Mutation(() => OrgUserSimDto)
    @UseMiddleware(Authentication)
    async updateUserRoaming(
        @Arg("data") data: UpdateUserServiceInput,
        @Ctx() ctx: Context
    ): Promise<OrgUserSimDto> {
        return this.userService.updateUserRoaming(data, parseHeaders(ctx));
    }
}
