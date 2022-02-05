import {
    Resolver,
    Query,
    UseMiddleware,
    PubSubEngine,
    PubSub,
    Arg,
} from "type-graphql";
import { Service } from "typedi";
import { NodeRFDto } from "../types";
import { NodeService } from "../service";
import { Authentication } from "../../../common/Authentication";
import { GRAPH_FILTER } from "../../../constants";

@Service()
@Resolver()
export class GetNodeRFKPIResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Query(() => [NodeRFDto])
    @UseMiddleware(Authentication)
    async getNodeRFKPI(
        @Arg("filter", () => GRAPH_FILTER) filter: GRAPH_FILTER,
        @PubSub() pubsub: PubSubEngine
    ): Promise<[NodeRFDto] | null> {
        const nodeRF = this.nodeService.nodeRF(filter);
        pubsub.publish("getNodeRFKPI", nodeRF);
        return nodeRF;
    }
}
