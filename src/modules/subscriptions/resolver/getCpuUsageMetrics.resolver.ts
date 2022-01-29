import { Resolver, Root, Subscription } from "type-graphql";
import { Service } from "typedi";
import { CpuUsageMetricsDto, CpuUsageMetricsResponse } from "../../node/types";

@Service()
@Resolver()
export class GetCpuUsageMetricsSubscriptionResolver {
    @Subscription(() => CpuUsageMetricsDto, {
        topics: "cpuUsageMetrics",
    })
    async getCpuUsageMetrics(
        @Root() data: CpuUsageMetricsResponse
    ): Promise<CpuUsageMetricsDto> {
        return data.data[0];
    }
}
