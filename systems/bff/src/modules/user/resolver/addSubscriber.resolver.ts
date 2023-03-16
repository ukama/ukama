import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseCookie } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { UserService } from "../service";
import { SubscriberDto, SubscriberInputDto } from "../types";

@Service()
@Resolver()
export class AddSubscriberResolver {
    constructor(private readonly userService: UserService) {}

    @Mutation(() => SubscriberDto)
    @UseMiddleware(Authentication)
    async addSubscriber(
        @Arg("data") data: SubscriberInputDto,
        @Ctx() ctx: Context,
    ): Promise<SubscriberDto> {
        return await this.userService.addSubscriber(data, parseCookie(ctx));
    }
}
