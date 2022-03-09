import {
    Arg,
    Ctx,
    Query,
    PubSub,
    Resolver,
    UseMiddleware,
    PubSubEngine,
} from "type-graphql";
import { Service } from "typedi";
import { MetricDto } from "../types";
import { NodeService } from "../service";
import { oneSecSleep } from "../../../utils";
import { getHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context, MetricsInputDTO } from "../../../common/types";

@Service()
@Resolver()
export class GetMetricsUptimeResolver {
    constructor(private readonly nodeService: NodeService) {}
    @Query(() => [MetricDto])
    @UseMiddleware(Authentication)
    async getMetricsUptime(
        @Ctx() ctx: Context,
        @Arg("data") data: MetricsInputDTO,
        @PubSub() pubsub: PubSubEngine
    ): Promise<MetricDto[] | null> {
        const metric = await this.nodeService.getSingleMetric(
            data,
            getHeaders(ctx),
            "uptime"
        );
        if (data.regPolling && metric && metric.length > 0) {
            for (const element of metric) {
                await oneSecSleep();
                pubsub.publish("metricUptime", [element]);
            }
        }
        return metric;
    }
}
