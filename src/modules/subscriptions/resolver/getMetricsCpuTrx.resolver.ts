import { Resolver, Root, Subscription } from "type-graphql";
import { Service } from "typedi";
import { MetricDto } from "../../node/types";

@Service()
@Resolver()
export class GetMetricsCpuTrxSubscriptionResolver {
    @Subscription(() => [MetricDto], {
        topics: "metricCpuTrx",
    })
    async getMetricsCpuTrx(@Root() data: [MetricDto]): Promise<MetricDto[]> {
        return data;
    }
}
