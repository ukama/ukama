import { Resolver, Root, Subscription } from "type-graphql";
import { Service } from "typedi";
import { NodePhysicalHealthDto } from "../../node/types";

@Service()
@Resolver()
export class GetNodesPhysicalHealthSubscriptionResolver {
    @Subscription(() => NodePhysicalHealthDto, {
        topics: "getNodeHealth",
    })
    async getNodePhysicalHealth(
        @Root() nodePhysicalHealth: NodePhysicalHealthDto
    ): Promise<NodePhysicalHealthDto> {
        return nodePhysicalHealth;
    }
}
