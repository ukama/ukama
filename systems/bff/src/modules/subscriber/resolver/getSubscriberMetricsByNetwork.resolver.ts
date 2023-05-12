import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { SubscriberService } from "../service";
import { SubscriberMetricsByNetworkDto } from "../types";

@Service()
@Resolver()
export class GetSubscriberMetricsByNetworkResolver {
    constructor(private readonly subService: SubscriberService) {}

    @Query(() => SubscriberMetricsByNetworkDto)
    @UseMiddleware(Authentication)
    async getSubscriberMetricsByNetwork(
        @Arg("networkId") networkId: string,
        @Ctx() ctx: Context
    ): Promise<SubscriberMetricsByNetworkDto> {
        return await this.subService.getSubMetricsByNetwork();
    }
}
