import {
    Resolver,
    Query,
    UseMiddleware,
    PubSubEngine,
    PubSub,
    Arg,
} from "type-graphql";
import { Service } from "typedi";
import { ThroughputMetricsDto } from "../types";
import { NodeService } from "../service";
import { Authentication } from "../../../common/Authentication";
import { GRAPH_FILTER } from "../../../constants";

@Service()
@Resolver()
export class GetThroughputMetricsResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Query(() => [ThroughputMetricsDto])
    @UseMiddleware(Authentication)
    async getThroughputMetrics(
        @Arg("filter", () => GRAPH_FILTER) filter: GRAPH_FILTER,
        @PubSub() pubsub: PubSubEngine
    ): Promise<[ThroughputMetricsDto]> {
        const throughputMetrics = this.nodeService.getThroughputMetrics(filter);
        pubsub.publish("throughputMetrics", throughputMetrics);
        return throughputMetrics;
    }
}
