import { Resolver, Root, Subscription } from "type-graphql";
import { Service } from "typedi";
import { MetricDto } from "../../node/types";

@Service()
@Resolver()
export class GetMetricsERABSubscriptionResolver {
    @Subscription(() => [MetricDto], {
        topics: "metricErab",
    })
    async getMetricsERAB(@Root() data: [MetricDto]): Promise<MetricDto[]> {
        return data;
    }
}
