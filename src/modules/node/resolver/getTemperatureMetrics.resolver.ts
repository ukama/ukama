import {
    Resolver,
    Arg,
    Query,
    UseMiddleware,
    PubSubEngine,
    PubSub,
} from "type-graphql";
import { Service } from "typedi";
import { NodeService } from "../service";
import { TemperatureMetricsResponse } from "../types";
import { PaginationDto } from "../../../common/types";
import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class GetTemperatureMetricsResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Query(() => TemperatureMetricsResponse)
    @UseMiddleware(Authentication)
    async getTemperatureMetrics(
        @Arg("data") data: PaginationDto,
        @PubSub() pubsub: PubSubEngine
    ): Promise<TemperatureMetricsResponse | null> {
        const temperatureMetrics = this.nodeService.temperatureMetrics(data);
        pubsub.publish("temperatureMetrics", temperatureMetrics);
        return temperatureMetrics;
    }
}
