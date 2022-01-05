import {
    Resolver,
    Query,
    UseMiddleware,
    PubSubEngine,
    PubSub,
} from "type-graphql";
import { Service } from "typedi";
import { NodeRFDto } from "../types";
import { NodeService } from "../service";

import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class GetNodeRFKPIResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Query(() => NodeRFDto)
    @UseMiddleware(Authentication)
    async getNodeRFKPI(@PubSub() pubsub: PubSubEngine): Promise<NodeRFDto> {
        const nodeRF = this.nodeService.nodeRF();
        pubsub.publish("getNodeRFKPI", nodeRF);
        return nodeRF;
    }
}
