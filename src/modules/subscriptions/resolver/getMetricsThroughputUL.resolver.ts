import { Resolver, Root, Subscription } from "type-graphql";
import { Service } from "typedi";
import { MetricDto } from "../../node/types";

@Service()
@Resolver()
export class GetMetricsThroughputULSubscriptionResolver {
    @Subscription(() => [MetricDto], {
        topics: "metricsThroughputUL",
    })
    async getMetricsThroughputUL(
        @Root() data: [MetricDto]
    ): Promise<MetricDto[]> {
        return data;
    }
}
