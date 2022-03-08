import { Resolver, Root, Subscription } from "type-graphql";
import { Service } from "typedi";
import { MetricDto } from "../../node/types";

@Service()
@Resolver()
export class GetMetricsMemoryTrxResolver {
    @Subscription(() => [MetricDto], {
        topics: "metricMemoryTrx",
    })
    async getMetricsMemoryTrx(@Root() data: [MetricDto]): Promise<MetricDto[]> {
        return data;
    }
}
