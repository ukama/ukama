import { Resolver, Arg, Ctx, Mutation, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { UserFistVisitInputDto, UserFistVisitResDto } from "../types";
import { UserService } from "../service";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { parseCookie } from "../../../common";

@Service()
@Resolver()
export class updateFirstVisitResolver {
    constructor(private readonly userService: UserService) {}

    @Mutation(() => UserFistVisitResDto)
    @UseMiddleware(Authentication)
    async updateFirstVisit(
        @Arg("data") data: UserFistVisitInputDto,
        @Ctx() ctx: Context,
    ): Promise<UserFistVisitResDto> {
        const user = await this.userService.updateFirstVisit(
            data,
            parseCookie(ctx),
        );
        await this.userService.updateFirstVisit(
            { firstVisit: data.firstVisit },
            parseCookie(ctx),
        );
        return user;
    }
}
