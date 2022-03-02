import { Service } from "typedi";
import { MetricDto } from "../types";
import { NodeService } from "../service";
import { getHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context, MetricsInputDTO } from "../../../common/types";
import { Resolver, Query, UseMiddleware, Arg, Ctx } from "type-graphql";

@Service()
@Resolver()
export class GetMetricsThroughputULResolver {
    constructor(private readonly nodeService: NodeService) {}
    @Query(() => [MetricDto])
    @UseMiddleware(Authentication)
    async getMetricsThroughputUL(
        @Ctx() ctx: Context,
        @Arg("data") data: MetricsInputDTO
    ): Promise<MetricDto[] | null> {
        return this.nodeService.metricsCpuTRX(data, getHeaders(ctx));
    }
}
