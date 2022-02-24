import {
    Resolver,
    Query,
    UseMiddleware,
    PubSubEngine,
    PubSub,
    Arg,
    Ctx,
} from "type-graphql";
import { Service } from "typedi";
import { NodeService } from "../service";
import { MetricsCpuTRXDto } from "../types";
import { getHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context, MetricsInputDTO } from "../../../common/types";

@Service()
@Resolver()
export class GetMetricsCpuTRXResolver {
    constructor(private readonly nodeService: NodeService) {}
    @Query(() => [MetricsCpuTRXDto])
    @UseMiddleware(Authentication)
    async getMetricsCpuTRX(
        @Ctx() ctx: Context,
        @Arg("data") data: MetricsInputDTO
    ): Promise<MetricsCpuTRXDto[] | null> {
        const metrics = this.nodeService.metricsCpuTRX(data, getHeaders(ctx));
        return metrics;
    }
}
