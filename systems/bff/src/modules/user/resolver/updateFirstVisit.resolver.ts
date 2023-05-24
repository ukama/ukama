import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { UserService } from "../service";
import { UserFistVisitInputDto, UserFistVisitResDto } from "../types";

@Service()
@Resolver()
export class updateFirstVisitResolver {
    constructor(private readonly userService: UserService) {}

    @Mutation(() => UserFistVisitResDto)
    @UseMiddleware(Authentication)
    async updateFirstVisit(
        @Arg("data") data: UserFistVisitInputDto,
        @Ctx() ctx: Context
    ): Promise<UserFistVisitResDto> {
        const user = await this.userService.updateFirstVisit(
            data,
            parseHeaders(ctx)
        );
        await this.userService.updateFirstVisit(
            { firstVisit: data.firstVisit },
            parseHeaders(ctx)
        );
        return user;
    }
}
