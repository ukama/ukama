import { Resolver, Root, Subscription } from "type-graphql";
import { Service } from "typedi";
import { TemperatureMetricsDto } from "../../node/types";

@Service()
@Resolver()
export class GetTemperatureMetricsSubscriptionResolver {
    @Subscription(() => TemperatureMetricsDto, {
        topics: "temperatureMetrics",
    })
    async getTemperatureMetrics(
        @Root() data: [TemperatureMetricsDto]
    ): Promise<TemperatureMetricsDto> {
        return data[data.length - 1];
    }
}
