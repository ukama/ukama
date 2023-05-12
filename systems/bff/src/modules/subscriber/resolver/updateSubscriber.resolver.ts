import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseCookie } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { BoolResponse, Context } from "../../../common/types";
import { SubscriberService } from "../service";
import { UpdateSubscriberInputDto } from "../types";

@Service()
@Resolver()
export class UpdateSubscriberResolver {
    constructor(private readonly userService: SubscriberService) {}

    @Mutation(() => BoolResponse)
    @UseMiddleware(Authentication)
    async updateSubscriber(
        @Arg("subscriberId") subscriberId: string,
        @Arg("data") data: UpdateSubscriberInputDto,
        @Ctx() ctx: Context
    ): Promise<BoolResponse> {
        return await this.userService.updateSubscriber(
            subscriberId,
            data,
            parseCookie(ctx)
        );
    }
}
