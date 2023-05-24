import {
    Arg,
    Ctx,
    PubSub,
    PubSubEngine,
    Query,
    Resolver,
    UseMiddleware,
} from "type-graphql";
import { Service } from "typedi";
import { Worker } from "worker_threads";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context, MetricsByTabInputDTO } from "../../../common/types";
import setupLogger from "../../../config/setupLogger";
import { getMetricsByTab } from "../../../utils";
import { NodeService } from "../service";
import { GetMetricsRes } from "../types";

const logger = setupLogger("Resolver");
const THREAD = "./MetricsThread.tsx";
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
            parseHeaders(ctx),
            metricsEndpoints
        );

        let next = false;
        if (data.regPolling) {
            const length = data.to - data.from;

            const workerData: any = { length, response };
            const worker = new Worker(THREAD, {
                workerData,
            });
            worker.on("message", (data: any) =>
                pubsub.publish("metricsByTab", data.metric)
            );
            worker.on("exit", (code: any) => {
                logger.info("Thread exited", code);
            });
        } else {
            for (const element of response) {
                if (!next && element.next) next = true;
            }
        }
        return {
            to: data.to,
            next: next,
            metrics: data.regPolling ? [] : response,
        };
    }
}
