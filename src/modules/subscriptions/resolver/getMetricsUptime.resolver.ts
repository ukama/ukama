import { Resolver, Root, Subscription } from "type-graphql";
import { Service } from "typedi";
import { MetricDto } from "../../node/types";

@Service()
@Resolver()
export class GetMetricsUptimeSubscriptionResolver {
    @Subscription(() => [MetricDto], {
        topics: "metricUptime",
    })
    async getMetricsUptime(@Root() data: [MetricDto]): Promise<MetricDto[]> {
        return data;
    }
}
