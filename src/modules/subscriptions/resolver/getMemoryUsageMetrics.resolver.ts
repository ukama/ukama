import { Service } from "typedi";
import { Resolver, Root, Subscription } from "type-graphql";
import {
    MemoryUsageMetricsDto,
    MemoryUsageMetricsResponse,
} from "../../node/types";

@Service()
@Resolver()
export class GetMemoryUsageMetricsSubscriptionResolver {
    @Subscription(() => MemoryUsageMetricsDto, {
        topics: "memoryUsageMetrics",
    })
    async getMemoryUsageMetrics(
        @Root() data: MemoryUsageMetricsResponse
    ): Promise<MemoryUsageMetricsDto> {
        return data.data[0];
    }
}
