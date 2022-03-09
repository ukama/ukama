import { Resolver, Root, Subscription } from "type-graphql";
import { Service } from "typedi";
import { MetricDto } from "../../node/types";

@Service()
@Resolver()
export class GetMetricsSubAttachedSubscriptionResolver {
    @Subscription(() => [MetricDto], {
        topics: "metricSubAttached",
    })
    async getMetricsSubAttached(
        @Root() data: [MetricDto]
    ): Promise<MetricDto[]> {
        return data;
    }
}
