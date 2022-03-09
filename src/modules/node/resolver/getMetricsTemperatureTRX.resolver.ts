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
export class GetMetricsTemperatureTRXResolver {
    constructor(private readonly nodeService: NodeService) {}
    @Query(() => [MetricDto])
    @UseMiddleware(Authentication)
    async getMetricsTemperatureTRX(
        @Ctx() ctx: Context,
        @Arg("data") data: MetricsInputDTO,
        @PubSub() pubsub: PubSubEngine
    ): Promise<MetricDto[] | null> {
        const metric = await this.nodeService.getSingleMetric(
            data,
            getHeaders(ctx),
            "temperaturetrx"
        );
        if (data.regPolling && metric && metric.length > 0) {
            for (const element of metric) {
                await oneSecSleep();
                pubsub.publish("metricTempTrx", [element]);
            }
        }
        return metric;
    }
}
