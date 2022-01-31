import { Resolver, Root, Subscription } from "type-graphql";
import { Service } from "typedi";
import { CpuUsageMetricsDto } from "../../node/types";

@Service()
@Resolver()
export class GetCpuUsageMetricsSubscriptionResolver {
    @Subscription(() => CpuUsageMetricsDto, {
        topics: "cpuUsageMetrics",
    })
    async getCpuUsageMetrics(
        @Root() data: [CpuUsageMetricsDto]
    ): Promise<CpuUsageMetricsDto> {
        return data[0];
    }
}
