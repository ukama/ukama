import { Resolver, Root, Subscription } from "type-graphql";
import { Service } from "typedi";
import { NodeRFDto } from "../../node/types";

@Service()
@Resolver()
export class GetNodeRFKPISubscriptionResolver {
    @Subscription(() => NodeRFDto, {
        topics: "getNodeRFKPI",
    })
    async getNodeRFKPI(@Root() nodeRF: NodeRFDto): Promise<NodeRFDto> {
        return nodeRF;
    }
}
