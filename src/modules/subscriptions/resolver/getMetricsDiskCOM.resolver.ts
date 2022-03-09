import { Resolver, Root, Subscription } from "type-graphql";
import { Service } from "typedi";
import { MetricDto } from "../../node/types";

@Service()
@Resolver()
export class GetMetricsDiskCOMSubscriptionResolver {
    @Subscription(() => [MetricDto], {
        topics: "metricDiskCom",
    })
    async getMetricsDiskCOM(@Root() data: [MetricDto]): Promise<MetricDto[]> {
        return data;
    }
}
