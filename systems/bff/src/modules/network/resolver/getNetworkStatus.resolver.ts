import {
    Resolver,
    Query,
    Arg,
    UseMiddleware,
    PubSubEngine,
    PubSub,
    Ctx,
} from "type-graphql";
import { Service } from "typedi";
import { NetworkDto } from "../types";
import { NetworkService } from "../service";
import { parseCookie } from "../../../common";
import { Context } from "../../../common/types";
import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class GetNetworkStatusResolver {
    constructor(private readonly networkService: NetworkService) {}

    @Query(() => NetworkDto)
    @UseMiddleware(Authentication)
    async getNetworkStatus(
        @Ctx() ctx: Context,
        @PubSub() pubsub: PubSubEngine,
    ): Promise<NetworkDto> {
        const network = this.networkService.getNetworkStatus(parseCookie(ctx));
        pubsub.publish("getNetworkStatus", network);
        return network;
    }
}
