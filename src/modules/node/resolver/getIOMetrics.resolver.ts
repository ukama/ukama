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
import { IOMetricsResponse } from "../types";
import { PaginationDto } from "../../../common/types";
import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class GetIOMetricsResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Query(() => IOMetricsResponse)
    @UseMiddleware(Authentication)
    async getIOMetrics(
        @Arg("data") data: PaginationDto,
        @PubSub() pubsub: PubSubEngine
    ): Promise<IOMetricsResponse | null> {
        const ioMetrics = this.nodeService.ioMetrics(data);
        pubsub.publish("ioMetrics", ioMetrics);
        return ioMetrics;
    }
}
