import { Resolver, Root, Subscription } from "type-graphql";
import { Service } from "typedi";
import { MetricDto } from "../../node/types";

@Service()
@Resolver()
export class GetMetricsCpuTRXSubscriptionResolver {
    @Subscription(() => [MetricDto], {
        topics: "metricCpuTrx",
    })
    async getMetricsCpuTRX(@Root() data: [MetricDto]): Promise<MetricDto[]> {
        return data;
    }
}
