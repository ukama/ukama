import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseCookie } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { BoolResponse, Context } from "../../../common/types";
import { SubscriberService } from "../service";

@Service()
@Resolver()
export class DeleteSubscriberResolver {
    constructor(private readonly userService: SubscriberService) {}

    @Mutation(() => BoolResponse)
    @UseMiddleware(Authentication)
    async deleteSubscriber(
        @Arg("subscriberId") subscriberId: string,
        @Ctx() ctx: Context
    ): Promise<BoolResponse> {
        return await this.userService.deleteSubscriber(
            subscriberId,
            parseCookie(ctx)
        );
    }
}
