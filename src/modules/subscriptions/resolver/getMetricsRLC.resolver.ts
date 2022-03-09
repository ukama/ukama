import { Resolver, Root, Subscription } from "type-graphql";
import { Service } from "typedi";
import { MetricDto } from "../../node/types";

@Service()
@Resolver()
export class GetMetricsRLCSubscriptionResolver {
    @Subscription(() => [MetricDto], {
        topics: "metricRlc",
    })
    async getMetricsRLC(@Root() data: [MetricDto]): Promise<MetricDto[]> {
        return data;
    }
}
