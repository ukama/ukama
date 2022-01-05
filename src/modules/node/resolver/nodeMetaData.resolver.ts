import {
    Resolver,
    Query,
    UseMiddleware,
    PubSub,
    PubSubEngine,
} from "type-graphql";
import { Service } from "typedi";
import { NodeMetaDataDto } from "../types";
import { NodeService } from "../service";

import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class GetNodeMetaDataResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Query(() => NodeMetaDataDto)
    @UseMiddleware(Authentication)
    async getNodeMetaData(
        @PubSub() pubsub: PubSubEngine
    ): Promise<NodeMetaDataDto> {
        const nodeMetaData = this.nodeService.nodeMetaData();
        pubsub.publish("getNodeMetaData", nodeMetaData);
        return nodeMetaData;
    }
}
