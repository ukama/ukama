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
import { NodeService } from "../service";
import { parseCookie } from "../../../common";
import { GetMetricsRes, MetricRes } from "../types";
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
            parseCookie(ctx),
            metricsEndpoints
        );

        let next = false;
        if (data.regPolling) {
            const length = data.to - data.from;
            for (let i = 0; i < length; i++) {
                const metric: MetricRes[] = [];
                for (const element of response) {
                    if (!next && element.next) next = true;
                    metric.push({
                        next: element.next,
                        type: element.type,
                        name: element.name,
                        data: element.data[i] ? [element.data[i]] : [],
                    });
                }
                await oneSecSleep();
                pubsub.publish("metricsByTab", metric);
            }
        } else {
            for (const element of response) {
                if (!next && element.next) next = true;
            }
        }

        return { to: data.to, next: next, metrics: response };
    }
}
