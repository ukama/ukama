import { Resolver, Root, Subscription } from "type-graphql";
import { Service } from "typedi";
import { MetricDto } from "../../node/types";

@Service()
@Resolver()
export class GetMetricsPowerSubscriptionResolver {
    @Subscription(() => [MetricDto], {
        topics: "metricPower",
    })
    async getMetricsPower(@Root() data: [MetricDto]): Promise<MetricDto[]> {
        return data;
    }
}
