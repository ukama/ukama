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
import { MetricDto } from "../types";
import { NodeService } from "../service";
import { getHeaders } from "../../../common";
import { oneSecSleep } from "../../../utils";
import { Authentication } from "../../../common/Authentication";
import { Context, MetricsInputDTO } from "../../../common/types";

@Service()
@Resolver()
export class GetCpuUsageMetricsResolver {
    constructor(private readonly cpuUsageMetrics: NodeService) {}

    @Query(() => [MetricDto])
    @UseMiddleware(Authentication)
    async getCpuUsageMetrics(
        @Ctx() ctx: Context,
        @Arg("data") data: MetricsInputDTO,
        @PubSub() pubsub: PubSubEngine
    ): Promise<MetricDto[] | null> {
        const cpuUsageMetrics = await this.cpuUsageMetrics.getSingleMetric(
            data,
            getHeaders(ctx),
            "cpu"
        );

        for (let i = 0; i < cpuUsageMetrics.length; i++) {
            await oneSecSleep();
            pubsub.publish("cpuUsageMetrics", [cpuUsageMetrics[i]]);
        }

        return cpuUsageMetrics;
    }
}
