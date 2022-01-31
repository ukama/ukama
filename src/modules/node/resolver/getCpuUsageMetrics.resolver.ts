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
import { CpuUsageMetricsResponse } from "../types";
import { PaginationDto } from "../../../common/types";
import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class GetCpuUsageMetricsResolver {
    constructor(private readonly cpuUsageMetrics: NodeService) {}

    @Query(() => CpuUsageMetricsResponse)
    @UseMiddleware(Authentication)
    async getCpuUsageMetrics(
        @Arg("data") data: PaginationDto,
        @PubSub() pubsub: PubSubEngine
    ): Promise<CpuUsageMetricsResponse | null> {
        const cpuUsageMetrics = this.cpuUsageMetrics.cpuUsageMetrics(data);
        pubsub.publish("cpuUsageMetrics", cpuUsageMetrics);
        return cpuUsageMetrics;
    }
}
