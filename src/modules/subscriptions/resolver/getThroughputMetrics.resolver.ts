import { Service } from "typedi";
import { ThroughputMetricsDto } from "../../node/types";
import { Resolver, Root, Subscription } from "type-graphql";

@Service()
@Resolver()
export class GetThroughputMetricsSubscriptionResolver {
    @Subscription(() => ThroughputMetricsDto, {
        topics: "throughputMetrics",
    })
    async getThroughputMetrics(
        @Root() data: [ThroughputMetricsDto]
    ): Promise<ThroughputMetricsDto> {
        return data[data.length - 1];
    }
}
