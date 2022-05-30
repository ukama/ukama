import { Resolver, Root, Subscription } from "type-graphql";
import { Service } from "typedi";
import { NetworkDto } from "../../network/types";

@Service()
@Resolver()
export class GetNetworkStatusSubscriptionResolver {
    @Subscription(() => NetworkDto, {
        topics: "getNetworkStatus",
    })
    async getNetworkStatus(@Root() network: NetworkDto): Promise<NetworkDto> {
        return network;
    }
}
