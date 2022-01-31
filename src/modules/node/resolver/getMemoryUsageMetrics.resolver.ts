import {
    Resolver,
    Arg,
    Query,
    UseMiddleware,
    PubSubEngine,
    PubSub,
} from "type-graphql";
import { Service } from "typedi";
import { NodeService } from "../service";
import { MemoryUsageMetricsResponse } from "../types";
import { PaginationDto } from "../../../common/types";
import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class GetMemoryUsageMetricsResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Query(() => MemoryUsageMetricsResponse)
    @UseMiddleware(Authentication)
    async getMemoryUsageMetrics(
        @Arg("data") data: PaginationDto,
        @PubSub() pubsub: PubSubEngine
    ): Promise<MemoryUsageMetricsResponse | null> {
        const ramUsageMetrics = this.nodeService.memoryUsageMetrics(data);
        pubsub.publish("memoryUsageMetrics", ramUsageMetrics);
        return ramUsageMetrics;
    }
}
