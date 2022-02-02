import {
    Resolver,
    Query,
    UseMiddleware,
    PubSubEngine,
    PubSub,
} from "type-graphql";
import { Service } from "typedi";
import { NodeService } from "../service";
import { IOMetricsDto } from "../types";
import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class GetIOMetricsResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Query(() => [IOMetricsDto])
    @UseMiddleware(Authentication)
    async getIOMetrics(
        @PubSub() pubsub: PubSubEngine
    ): Promise<[IOMetricsDto] | null> {
        const ioMetrics = this.nodeService.ioMetrics();
        pubsub.publish("ioMetrics", ioMetrics);
        return ioMetrics;
    }
}
