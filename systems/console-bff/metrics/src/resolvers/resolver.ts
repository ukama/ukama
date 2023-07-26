import {
  Arg,
  Args,
  PubSub,
  PubSubEngine,
  Query,
  Resolver,
  Root,
  Subscription,
} from "type-graphql";
import { Worker } from "worker_threads";

import { METRIC_API_GW_SOCKET, STORAGE_KEY } from "../../../common/configs";
import { logger } from "../../../common/logger";
import { storeInStorage } from "../../../common/storage";
import { getTimestampCount } from "../../../common/utils";
import { getLatestMetric, getMetricRange } from "../datasource/metrics-api";
import {
  GetLatestMetricInput,
  GetMetricRangeInput,
  LatestMetricRes,
  MetricRes,
  SubMetricRangeInput,
} from "./types";

const WS_THREAD = "./threads/MetricsWSThread.ts";

@Resolver(MetricRes)
class MetricResolvers {
  @Query(() => LatestMetricRes)
  async getLatestMetric(@Arg("data") data: GetLatestMetricInput) {
    return await getLatestMetric(data);
  }

  @Query(() => MetricRes)
  async getMetricRange(
    @Arg("data") data: GetMetricRangeInput,
    @PubSub() pubSub: PubSubEngine
  ) {
    const { type, orgId, userId, nodeId, withSubscription } = data;

    if (withSubscription) {
      const workerData: any = {
        type,
        orgId,
        userId,
        url: `${METRIC_API_GW_SOCKET}/v1/live/metric?interval=1&metric=${type}`,
        key: STORAGE_KEY,
      };
      const worker = new Worker(WS_THREAD, {
        workerData,
      });
      worker.on("message", (_data: any) => {
        if (!_data.isError) {
          const res = JSON.parse(_data.data);
          const result = res.data.result[0];
          pubSub.publish(`metric-${type}`, {
            env: result.metric.env,
            nodeid: nodeId,
            type: type,
            value: result.value,
          } as LatestMetricRes);
        }
      });
      worker.on("exit", (code: any) => {
        logger.info(
          `WS_THREAD exited with code [${code}] for ${orgId}/${userId}/${type}`
        );
      });
    }

    return await getMetricRange(data);
  }

  @Subscription(() => LatestMetricRes, {
    topics: ({ args }) => `metric-${args.type}`,
    filter: ({ payload, args }) => {
      return args.nodeId === payload.nodeid;
    },
  })
  async getMetric(
    @Root() payload: LatestMetricRes,
    @Args() args: SubMetricRangeInput
  ): Promise<LatestMetricRes> {
    await storeInStorage(
      `${args.orgId}/${args.userId}/${args.type}`,
      getTimestampCount("0")
    );
    return payload;
  }
}

export default MetricResolvers;
