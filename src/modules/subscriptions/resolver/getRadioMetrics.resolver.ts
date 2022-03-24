import { Service } from "typedi";
import { MetricRes } from "../../node/types";
import { Resolver, Root, Subscription } from "type-graphql";

@Service()
@Resolver()
export class GetRadioMetricsSubscriptionResolver {
    @Subscription(() => [MetricRes], {
        topics: "radioMetrics",
    })
    async getRadioMetrics(@Root() data: [MetricRes]): Promise<MetricRes[]> {
        return data;
    }
}
