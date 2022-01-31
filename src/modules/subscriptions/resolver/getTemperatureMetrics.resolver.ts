import { Resolver, Root, Subscription } from "type-graphql";
import { Service } from "typedi";
import {
    TemperatureMetricsDto,
    TemperatureMetricsResponse,
} from "../../node/types";

@Service()
@Resolver()
export class GetTemperatureMetricsSubscriptionResolver {
    @Subscription(() => TemperatureMetricsDto, {
        topics: "temperatureMetrics",
    })
    async getTemperatureMetrics(
        @Root() data: TemperatureMetricsResponse
    ): Promise<TemperatureMetricsDto> {
        return data.data[0];
    }
}
