import { Resolver, Root, Subscription } from "type-graphql";
import { Service } from "typedi";
import { NodeRFDto, NodeRFDtoResponse } from "../../node/types";

@Service()
@Resolver()
export class GetNodeRFKPISubscriptionResolver {
    @Subscription(() => NodeRFDto, {
        topics: "getNodeRFKPI",
    })
    async getNodeRFKPI(@Root() data: NodeRFDtoResponse): Promise<NodeRFDto> {
        return data.data[0];
    }
}
