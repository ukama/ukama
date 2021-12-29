import {
    Resolver,
    Query,
    UseMiddleware,
    PubSubEngine,
    PubSub,
} from "type-graphql";
import { Service } from "typedi";
import { NodePhysicalHealthDto } from "../types";
import { NodeService } from "../service";

import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class GetNodePhysicalHealthResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Query(() => NodePhysicalHealthDto)
    @UseMiddleware(Authentication)
    async getNodePhysicalHealth(
        @PubSub() pubsub: PubSubEngine
    ): Promise<NodePhysicalHealthDto> {
        const nodeHealth = this.nodeService.nodePhysicalHealth();
        pubsub.publish("getNodeHealth", nodeHealth);
        return nodeHealth;
    }
}
