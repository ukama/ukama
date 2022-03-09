import { Resolver, Root, Subscription } from "type-graphql";
import { Service } from "typedi";
import { MetricDto } from "../../node/types";

@Service()
@Resolver()
export class GetMetricsMemoryCOMSubscriptionResolver {
    @Subscription(() => [MetricDto], {
        topics: "metricMemoryCom",
    })
    async getMetricsMemoryCOM(@Root() data: [MetricDto]): Promise<MetricDto[]> {
        return data;
    }
}
