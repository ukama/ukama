import { Service } from "typedi";
import { MemoryUsageMetricsDto } from "../../node/types";
import { Resolver, Root, Subscription } from "type-graphql";

@Service()
@Resolver()
export class GetMemoryUsageMetricsSubscriptionResolver {
    @Subscription(() => MemoryUsageMetricsDto, {
        topics: "memoryUsageMetrics",
    })
    async getMemoryUsageMetrics(
        @Root() data: [MemoryUsageMetricsDto]
    ): Promise<MemoryUsageMetricsDto> {
        return data[data.length - 1];
    }
}
