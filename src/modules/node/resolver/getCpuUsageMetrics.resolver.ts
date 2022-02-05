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
import { CpuUsageMetricsDto } from "../types";
import { Authentication } from "../../../common/Authentication";
import { GRAPH_FILTER } from "../../../constants";

@Service()
@Resolver()
export class GetCpuUsageMetricsResolver {
    constructor(private readonly cpuUsageMetrics: NodeService) {}

    @Query(() => [CpuUsageMetricsDto])
    @UseMiddleware(Authentication)
    async getCpuUsageMetrics(
        @Arg("filter", () => GRAPH_FILTER) filter: GRAPH_FILTER,
        @PubSub() pubsub: PubSubEngine
    ): Promise<[CpuUsageMetricsDto] | null> {
        const cpuUsageMetrics = this.cpuUsageMetrics.cpuUsageMetrics(filter);
        pubsub.publish("cpuUsageMetrics", cpuUsageMetrics);
        return cpuUsageMetrics;
    }
}
