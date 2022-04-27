import { Service } from "typedi";
import { MetricRes } from "../../node/types";
import { Resolver, Root, Subscription } from "type-graphql";

@Service()
@Resolver()
export class GetMetricsByTabSubscriptionResolver {
    @Subscription(() => [MetricRes], {
        topics: "metricsByTab",
    })
    async getMetricsByTab(@Root() data: [MetricRes]): Promise<MetricRes[]> {
        return data;
    }
}
