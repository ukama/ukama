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
import { RamUsageMetricsResponse } from "../types";
import { PaginationDto } from "../../../common/types";
import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class GetRamUsageMetricsResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Query(() => RamUsageMetricsResponse)
    @UseMiddleware(Authentication)
    async getRamUsageMetrics(
        @Arg("data") data: PaginationDto,
        @PubSub() pubsub: PubSubEngine
    ): Promise<RamUsageMetricsResponse | null> {
        const ramUsageMetrics = this.nodeService.ramUsageMetrics(data);
        pubsub.publish("ramUsageMetrics", ramUsageMetrics);
        return ramUsageMetrics;
    }
}
