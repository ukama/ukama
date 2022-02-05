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
import { IOMetricsDto } from "../types";
import { Authentication } from "../../../common/Authentication";
import { GRAPH_FILTER } from "../../../constants";

@Service()
@Resolver()
export class GetIOMetricsResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Query(() => [IOMetricsDto])
    @UseMiddleware(Authentication)
    async getIOMetrics(
        @Arg("filter", () => GRAPH_FILTER) filter: GRAPH_FILTER,
        @PubSub() pubsub: PubSubEngine
    ): Promise<[IOMetricsDto] | null> {
        const ioMetrics = this.nodeService.ioMetrics(filter);
        pubsub.publish("ioMetrics", ioMetrics);
        return ioMetrics;
    }
}
