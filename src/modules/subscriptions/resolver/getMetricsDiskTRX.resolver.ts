import { Resolver, Root, Subscription } from "type-graphql";
import { Service } from "typedi";
import { MetricDto } from "../../node/types";

@Service()
@Resolver()
export class GetMetricsDiskTRXSubscriptionResolver {
    @Subscription(() => [MetricDto], {
        topics: "metricDiskTrx",
    })
    async getMetricsDiskTRX(@Root() data: [MetricDto]): Promise<MetricDto[]> {
        return data;
    }
}
