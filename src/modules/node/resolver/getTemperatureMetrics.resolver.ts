import {
    Resolver,
    Query,
    UseMiddleware,
    PubSubEngine,
    PubSub,
    Arg,
} from "type-graphql";
import { Service } from "typedi";
import { NodeService } from "../service";
import { TemperatureMetricsDto } from "../types";
import { Authentication } from "../../../common/Authentication";
import { GRAPH_FILTER } from "../../../constants";

@Service()
@Resolver()
export class GetTemperatureMetricsResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Query(() => [TemperatureMetricsDto])
    @UseMiddleware(Authentication)
    async getTemperatureMetrics(
        @Arg("filter", () => GRAPH_FILTER) filter: GRAPH_FILTER,
        @PubSub() pubsub: PubSubEngine
    ): Promise<[TemperatureMetricsDto] | null> {
        const temperatureMetrics = this.nodeService.temperatureMetrics(filter);
        pubsub.publish("temperatureMetrics", temperatureMetrics);
        return temperatureMetrics;
    }
}
