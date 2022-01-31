import {
    Resolver,
    Query,
    UseMiddleware,
    PubSubEngine,
    PubSub,
} from "type-graphql";
import { Service } from "typedi";
import { ThroughputMetricsDto } from "../types";
import { NodeService } from "../service";
import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class GetThroughputMetricsResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Query(() => [ThroughputMetricsDto])
    @UseMiddleware(Authentication)
    async getThroughputMetrics(
        @PubSub() pubsub: PubSubEngine
    ): Promise<[ThroughputMetricsDto]> {
        const throughputMetrics = this.nodeService.getThroughputMetrics();
        pubsub.publish("throughputMetrics", throughputMetrics);
        return throughputMetrics;
    }
}
