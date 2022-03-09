import { Resolver, Root, Subscription } from "type-graphql";
import { Service } from "typedi";
import { MetricDto } from "../../node/types";

@Service()
@Resolver()
export class GetMetricsTxPowerSubscriptionResolver {
    @Subscription(() => [MetricDto], {
        topics: "metricTxPower",
    })
    async getMetricsTxPower(@Root() data: [MetricDto]): Promise<MetricDto[]> {
        return data;
    }
}
