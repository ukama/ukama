import { Resolver, Root, Subscription } from "type-graphql";
import { Service } from "typedi";
import { MetricDto } from "../../node/types";

@Service()
@Resolver()
export class GetMetricsThroughputDLSubscriptionResolver {
    @Subscription(() => [MetricDto], {
        topics: "metricsThroughputDL",
    })
    async getMetricsThroughputDL(
        @Root() data: [MetricDto]
    ): Promise<MetricDto[]> {
        return data;
    }
}
