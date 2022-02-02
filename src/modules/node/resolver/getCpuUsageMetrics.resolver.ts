import {
    Resolver,
    Query,
    UseMiddleware,
    PubSubEngine,
    PubSub,
} from "type-graphql";
import { Service } from "typedi";
import { NodeService } from "../service";
import { CpuUsageMetricsDto } from "../types";
import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class GetCpuUsageMetricsResolver {
    constructor(private readonly cpuUsageMetrics: NodeService) {}

    @Query(() => [CpuUsageMetricsDto])
    @UseMiddleware(Authentication)
    async getCpuUsageMetrics(
        @PubSub() pubsub: PubSubEngine
    ): Promise<[CpuUsageMetricsDto] | null> {
        const cpuUsageMetrics = this.cpuUsageMetrics.cpuUsageMetrics();
        pubsub.publish("cpuUsageMetrics", cpuUsageMetrics);
        return cpuUsageMetrics;
    }
}
