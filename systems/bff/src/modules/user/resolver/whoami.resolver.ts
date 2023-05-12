import { Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseCookie } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { UserService } from "../service";
import { WhoamiDto } from "../types";

@Service()
@Resolver()
export class WhoamiResolver {
    constructor(private readonly userService: UserService) {}
    @Query(() => WhoamiDto)
    @UseMiddleware(Authentication)
    async whoami(@Ctx() ctx: Context): Promise<WhoamiDto> {
        return this.userService.whoami(parseCookie(ctx));
    }
}
