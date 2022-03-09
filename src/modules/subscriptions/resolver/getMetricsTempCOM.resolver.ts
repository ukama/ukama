import { Resolver, Root, Subscription } from "type-graphql";
import { Service } from "typedi";
import { MetricDto } from "../../node/types";

@Service()
@Resolver()
export class GetMetricsTempCOMSubscriptionResolver {
    @Subscription(() => [MetricDto], {
        topics: "metricTempCom",
    })
    async getMetricsTempCOM(@Root() data: [MetricDto]): Promise<MetricDto[]> {
        return data;
    }
}
