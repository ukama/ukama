import { Resolver, Root, Subscription } from "type-graphql";
import { Service } from "typedi";
import { MetricDto } from "../../node/types";

@Service()
@Resolver()
export class GetMetricsMemoryTRXSubscriptionResolver {
    @Subscription(() => [MetricDto], {
        topics: "metricMemoryTrx",
    })
    async getMetricsMemoryTRX(@Root() data: [MetricDto]): Promise<MetricDto[]> {
        return data;
    }
}
