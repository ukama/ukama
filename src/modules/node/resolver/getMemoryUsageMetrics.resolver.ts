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
import { MemoryUsageMetricsDto } from "../types";
import { Authentication } from "../../../common/Authentication";
import { GRAPH_FILTER } from "../../../constants";

@Service()
@Resolver()
export class GetMemoryUsageMetricsResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Query(() => [MemoryUsageMetricsDto])
    @UseMiddleware(Authentication)
    async getMemoryUsageMetrics(
        @Arg("filter", () => GRAPH_FILTER) filter: GRAPH_FILTER,
        @PubSub() pubsub: PubSubEngine
    ): Promise<[MemoryUsageMetricsDto] | null> {
        const memoryUsageMetrics = this.nodeService.memoryUsageMetrics(filter);
        pubsub.publish("memoryUsageMetrics", memoryUsageMetrics);
        return memoryUsageMetrics;
    }
}
