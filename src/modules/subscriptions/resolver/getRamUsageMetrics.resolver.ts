import { Service } from "typedi";
import { Resolver, Root, Subscription } from "type-graphql";
import { RamUsageMetricsDto, RamUsageMetricsResponse } from "../../node/types";

@Service()
@Resolver()
export class GetRamUsageMetricsSubscriptionResolver {
    @Subscription(() => RamUsageMetricsDto, {
        topics: "ramUsageMetrics",
    })
    async getRamUsageMetrics(
        @Root() data: RamUsageMetricsResponse
    ): Promise<RamUsageMetricsDto> {
        return data.data[0];
    }
}
