import { Resolver, Root, Subscription } from "type-graphql";
import { Service } from "typedi";
import { MetricDto } from "../../node/types";

@Service()
@Resolver()
export class GetMetricsTempTRXSubscriptionResolver {
    @Subscription(() => [MetricDto], {
        topics: "metricTempTrx",
    })
    async getMetricsTempTRX(@Root() data: [MetricDto]): Promise<MetricDto[]> {
        return data;
    }
}
