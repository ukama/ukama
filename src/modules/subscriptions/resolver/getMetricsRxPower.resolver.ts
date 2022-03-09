import { Resolver, Root, Subscription } from "type-graphql";
import { Service } from "typedi";
import { MetricDto } from "../../node/types";

@Service()
@Resolver()
export class GetMetricsRxPowerSubscriptionResolver {
    @Subscription(() => [MetricDto], {
        topics: "metricRxPower",
    })
    async getMetricsRxPower(@Root() data: [MetricDto]): Promise<MetricDto[]> {
        return data;
    }
}
