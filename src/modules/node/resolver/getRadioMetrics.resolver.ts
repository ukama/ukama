import { Service } from "typedi";
import { MetricRes } from "../types";
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

const METRICS_LIST = [
    {
        key: 0,
        id: "txMetric",
        title: "TX Power",
    },
    {
        key: 1,
        id: "rxMetric",
        title: "RX Power",
    },
    {
        key: 2,
        id: "paMetric",
        title: "PA Power",
    },
];

@Service()
@Resolver()
export class GetRadioMetricsResolver {
    constructor(private readonly nodeService: NodeService) {}
    @Query(() => [MetricRes])
    @UseMiddleware(Authentication)
    async getRadioMetrics(
        @Ctx() ctx: Context,
        @PubSub() pubsub: PubSubEngine,
        @Arg("data") data: MetricsInputDTO
    ): Promise<MetricRes[] | null> {
        const rxMetric = await this.nodeService.getSingleMetric(
            data,
            getHeaders(ctx),
            "rxpower"
        );
        const paMetric = await this.nodeService.getSingleMetric(
            data,
            getHeaders(ctx),
            "papower"
        );
        const txMetric = await this.nodeService.getSingleMetric(
            data,
            getHeaders(ctx),
            "txpower"
        );

        const metricRes: MetricRes[] = [];

        if (
            rxMetric &&
            paMetric &&
            txMetric &&
            data.regPolling &&
            rxMetric.length > 0 &&
            txMetric.length > 0 &&
            paMetric.length > 0
        ) {
            for (let i = 0; i < rxMetric.length; i++) {
                const metric: MetricRes[] = [];
                metric.push({
                    type: "",
                    name: METRICS_LIST[0].title,
                    data: [txMetric[i]],
                });
                metric.push({
                    type: "",
                    name: METRICS_LIST[1].title,
                    data: [rxMetric[i]],
                });
                metric.push({
                    type: "",
                    name: METRICS_LIST[2].title,
                    data: [paMetric[i]],
                });

                await oneSecSleep();
                pubsub.publish("radioMetrics", metric);
            }
        }
        metricRes.push({
            type: "",
            name: METRICS_LIST[0].title,
            data: txMetric,
        });
        metricRes.push({
            type: "",
            name: METRICS_LIST[1].title,
            data: rxMetric,
        });
        metricRes.push({
            type: "",
            name: METRICS_LIST[2].title,
            data: paMetric,
        });
        return metricRes;
    }
}
