import { Resolver, Root, Subscription } from "type-graphql";
import { Service } from "typedi";
import { MetricDto } from "../../node/types";

@Service()
@Resolver()
export class GetMetricsRRCSubscriptionResolver {
    @Subscription(() => [MetricDto], {
        topics: "metricRrc",
    })
    async getMetricsRRC(@Root() data: [MetricDto]): Promise<MetricDto[]> {
        return data;
    }
}
