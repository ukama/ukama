import { Service } from "typedi";
import { MetricDto } from "../types";
import { NodeService } from "../service";
import { getHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context, MetricsInputDTO } from "../../../common/types";
import {
    Resolver,
    Query,
    UseMiddleware,
    Arg,
    Ctx,
    PubSubEngine,
    PubSub,
} from "type-graphql";
import { oneSecSleep } from "../../../utils";

@Service()
@Resolver()
export class GetMetricsRRCResolver {
    constructor(private readonly nodeService: NodeService) {}
    @Query(() => [MetricDto])
    @UseMiddleware(Authentication)
    async getMetricsRRC(
        @Ctx() ctx: Context,
        @PubSub() pubsub: PubSubEngine,
        @Arg("data") data: MetricsInputDTO
    ): Promise<MetricDto[] | null> {
        const metric = await this.nodeService.getSingleMetric(
            data,
            getHeaders(ctx),
            "rrc"
        );
        if (data.regPolling && metric && metric.length > 0) {
            for (const element of metric) {
                await oneSecSleep();
                pubsub.publish("metricRrc", [element]);
            }
        }
        return metric;
    }
}
