import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseCookie } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { SubscriberService } from "../service";
import { SubscriberDto } from "../types";

@Service()
@Resolver()
export class GetSubscriberResolver {
    constructor(private readonly userService: SubscriberService) {}

    @Query(() => SubscriberDto)
    @UseMiddleware(Authentication)
    async getSubscriber(
        @Arg("subscriberId") subscriberId: string,
        @Ctx() ctx: Context
    ): Promise<SubscriberDto> {
        return await this.userService.getSubscriber(
            subscriberId,
            parseCookie(ctx)
        );
    }
}
