import { Resolver, Root, Subscription } from "type-graphql";
import { Service } from "typedi";
import { MetricDto } from "../../node/types";

@Service()
@Resolver()
export class GetMetricsSubActiveSubscriptionResolver {
    @Subscription(() => [MetricDto], {
        topics: "metricSubActive",
    })
    async getMetricsSubActive(@Root() data: [MetricDto]): Promise<MetricDto[]> {
        return data;
    }
}
