import {
    Resolver,
    Query,
    UseMiddleware,
    PubSubEngine,
    PubSub,
} from "type-graphql";
import { Service } from "typedi";
import { NodeService } from "../service";
import { TemperatureMetricsDto } from "../types";
import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class GetTemperatureMetricsResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Query(() => [TemperatureMetricsDto])
    @UseMiddleware(Authentication)
    async getTemperatureMetrics(
        @PubSub() pubsub: PubSubEngine
    ): Promise<[TemperatureMetricsDto] | null> {
        const temperatureMetrics = this.nodeService.temperatureMetrics();
        pubsub.publish("temperatureMetrics", temperatureMetrics);
        return temperatureMetrics;
    }
}
