import { Resolver, Root, Subscription } from "type-graphql";
import { Service } from "typedi";
import { MetricDto } from "../../node/types";

@Service()
@Resolver()
export class GetMetricsCpuCOMSubscriptionResolver {
    @Subscription(() => [MetricDto], {
        topics: "metricCpuCom",
    })
    async getMetricsCpuCOM(@Root() data: [MetricDto]): Promise<MetricDto[]> {
        return data;
    }
}
