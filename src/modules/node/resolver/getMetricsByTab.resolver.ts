import {
    Arg,
    Ctx,
    Query,
    PubSub,
    Resolver,
    PubSubEngine,
    UseMiddleware,
} from "type-graphql";
import { Service } from "typedi";
import { GetMetricsRes, MetricRes } from "../types";
import { NodeService } from "../service";
import { getHeaders } from "../../../common";
import { getMetricsByTab, oneSecSleep } from "../../../utils";
import { Authentication } from "../../../common/Authentication";
import { Context, MetricsByTabInputDTO } from "../../../common/types";

@Service()
@Resolver()
export class GetMetricsByTabResolver {
    constructor(private readonly nodeService: NodeService) {}
    @Query(() => GetMetricsRes)
    @UseMiddleware(Authentication)
    async getMetricsByTab(
        @Ctx() ctx: Context,
        @PubSub() pubsub: PubSubEngine,
        @Arg("data") data: MetricsByTabInputDTO
    ): Promise<GetMetricsRes | null> {
        const metricsEndpoints = getMetricsByTab(data.nodeType, data.tab);
        const response = await this.nodeService.getMultipleMetrics(
            data,
            getHeaders(ctx),
            metricsEndpoints
        );

        if (data.regPolling) {
            const length = data.to - data.from;
            for (let i = 0; i < length; i++) {
                const metric: MetricRes[] = [];
                for (const element of response) {
                    metric.push({
                        type: element.type,
                        name: element.name,
                        data: element.data[i] ? [element.data[i]] : [],
                    });
                }
                await oneSecSleep();
                pubsub.publish("metricsByTab", metric);
            }
        }

        return { to: data.to, metrics: response };
    }
}
