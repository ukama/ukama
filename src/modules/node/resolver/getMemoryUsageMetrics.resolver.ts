import {
    Resolver,
    Query,
    UseMiddleware,
    PubSubEngine,
    PubSub,
} from "type-graphql";
import { Service } from "typedi";
import { NodeService } from "../service";
import { MemoryUsageMetricsDto } from "../types";
import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class GetMemoryUsageMetricsResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Query(() => [MemoryUsageMetricsDto])
    @UseMiddleware(Authentication)
    async getMemoryUsageMetrics(
        @PubSub() pubsub: PubSubEngine
    ): Promise<[MemoryUsageMetricsDto] | null> {
        const memoryUsageMetrics = this.nodeService.memoryUsageMetrics();
        pubsub.publish("memoryUsageMetrics", memoryUsageMetrics);
        return memoryUsageMetrics;
    }
}
