import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { SubscriberService } from "../service";
import { SubscribersResDto } from "../types";

@Service()
@Resolver()
export class GetSubscribersByNetworkResolver {
    constructor(private readonly userService: SubscriberService) {}

    @Query(() => SubscribersResDto)
    @UseMiddleware(Authentication)
    async getSubscribersByNetwork(
        @Arg("networkId") networkId: string,
        @Ctx() ctx: Context
    ): Promise<SubscribersResDto> {
        return await this.userService.getSubscribersByNetwork(
            networkId,
            parseHeaders(ctx)
        );
    }
}
