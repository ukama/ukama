import {
    Resolver,
    Query,
    UseMiddleware,
    PubSubEngine,
    PubSub,
    Arg,
} from "type-graphql";
import { Service } from "typedi";
import { NodeService } from "../service";
import { NodeRFDtoResponse } from "../types";
import { PaginationDto } from "../../../common/types";
import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class GetNodeRFKPIResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Query(() => NodeRFDtoResponse)
    @UseMiddleware(Authentication)
    async getNodeRFKPI(
        @Arg("data") data: PaginationDto,
        @PubSub() pubsub: PubSubEngine
    ): Promise<NodeRFDtoResponse | null> {
        const nodeRF = this.nodeService.nodeRF(data);
        pubsub.publish("getNodeRFKPI", nodeRF);
        return nodeRF;
    }
}
