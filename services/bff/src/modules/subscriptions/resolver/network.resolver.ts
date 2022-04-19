import { Resolver, Root, Subscription } from "type-graphql";
import { Service } from "typedi";
import { NetworkDto } from "../../network/types";

@Service()
@Resolver()
export class NetworkSubscriptionResolver {
    @Subscription(() => NetworkDto, {
        topics: "getNetwork",
    })
    async getNetwork(@Root() network: NetworkDto): Promise<NetworkDto> {
        return network;
    }
}
