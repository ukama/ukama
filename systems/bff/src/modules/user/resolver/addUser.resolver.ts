import { Resolver, Arg, Ctx, Mutation, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { UserInputDto, UserResDto } from "../types";
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
        @Arg("data") data: UserInputDto,
        @Ctx() ctx: Context
    ): Promise<UserResDto> {
        const user = await this.userService.addUser(data, parseCookie(ctx));
        await this.userService.updateUserRoaming(
            { simId: user?.iccid || "", userId: user.id, status: data.status },
            parseCookie(ctx)
        );
        return user;
    }
}
