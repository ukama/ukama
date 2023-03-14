import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseCookie } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { BoolResponse, Context } from "../../../common/types";
import { UserService } from "../service";

@Service()
@Resolver()
export class DeleteUserResolver {
    constructor(private readonly userService: UserService) {}

    @Mutation(() => BoolResponse)
    @UseMiddleware(Authentication)
    async deleteUser(
        @Arg("uuid") uuid: string,
        @Ctx() ctx: Context,
    ): Promise<BoolResponse> {
        return this.userService.deleteUser(uuid, parseCookie(ctx));
    }
}
