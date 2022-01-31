import { Resolver, Root, Subscription } from "type-graphql";
import { Service } from "typedi";
import { NodeRFDto } from "../../node/types";

@Service()
@Resolver()
export class GetNodeRFKPISubscriptionResolver {
    @Subscription(() => NodeRFDto, {
        topics: "getNodeRFKPI",
    })
    async getNodeRFKPI(@Root() data: [NodeRFDto]): Promise<NodeRFDto> {
        return data[0];
    }
}
