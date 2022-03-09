import { Resolver, Root, Subscription } from "type-graphql";
import { Service } from "typedi";
import { MetricDto } from "../../node/types";

@Service()
@Resolver()
export class GetMetricsPaPowerSubscriptionResolver {
    @Subscription(() => [MetricDto], {
        topics: "metricPaPower",
    })
    async getMetricsPaPower(@Root() data: [MetricDto]): Promise<MetricDto[]> {
        return data;
    }
}
