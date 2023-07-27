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

import { METRIC_API_GW_SOCKET, STORAGE_KEY } from "../../common/configs";
import { logger } from "../../common/logger";
import { removeKeyFromStorage, storeInStorage } from "../../common/storage";
import { getTimestampCount } from "../../common/utils";
import {
  getLatestMetric,
  getMetricRange,
  getNodeRangeMetric,
} from "../datasource/metrics-api";
import {
  GetLatestMetricInput,
  GetMetricRangeInput,
  LatestMetricRes,
  MetricRes,
  SubMetricRangeInput,
} from "./types";

const WS_THREAD = "./threads/MetricsWSThread.ts";

const getErrorRes = (msg: string) =>
  ({
    env: "",
    msg: msg,
    type: "",
    nodeid: "",
    values: [],
    success: false,
  } as MetricRes);

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
    const { type, orgId, userId, nodeId, withSubscription, from } = data;
    if (from === 0) throw new Error("Argument 'from' can't be zero.");
    const res = await getMetricRange(data);
    if (withSubscription && res.env) {
      const workerData: any = {
        type,
        orgId,
        userId,
        url: `${METRIC_API_GW_SOCKET}/v1/live/metric?interval=1&metric=${type}`,
        key: STORAGE_KEY,
        timestamp: from,
      };
      const worker = new Worker(WS_THREAD, {
        workerData,
      });
      worker.on("message", (_data: any) => {
        if (!_data.isError) {
          const res = JSON.parse(_data.data);
          const result = res.data.result[0];
          if (result && result.metric && result.value.length > 0) {
            pubSub.publish(`metric-${type}`, {
              success: true,
              msg: "success",
              env: result.metric.env,
              nodeid: nodeId,
              type: type,
              value: result.value,
            } as LatestMetricRes);
          } else {
            return getErrorRes("No metric data found");
          }
        }
      });
      worker.on("exit", (code: any) => {
        removeKeyFromStorage(`${orgId}/${userId}/${type}/${from}`);
        logger.info(
          `WS_THREAD exited with code [${code}] for ${orgId}/${userId}/${type}`
        );
      });
    }

    return res;
  }

  @Query(() => MetricRes)
  async getNodeRangeMetric(
    @Arg("data") data: GetMetricRangeInput,
    @PubSub() pubSub: PubSubEngine
  ) {
    const { type, orgId, userId, nodeId, withSubscription, from } = data;
    if (from === 0) throw new Error("Argument 'from' can't be zero.");
    const res = await getNodeRangeMetric(data);
    if (withSubscription && res.env) {
      const workerData: any = {
        type,
        orgId,
        userId,
        timestamp: from,
        key: STORAGE_KEY,
        url: `${METRIC_API_GW_SOCKET}/v1/live/metric?interval=1&metric=${type}`,
      };
      const worker = new Worker(WS_THREAD, {
        workerData,
      });
      worker.on("message", (_data: any) => {
        if (!_data.isError) {
          const res = JSON.parse(_data.data);
          const result = res.data.result[0];
          if (result && result.metric && result.value.length > 0) {
            pubSub.publish(`metric-${type}`, {
              env: result.metric.env,
              nodeid: nodeId,
              type: type,
              value:
                result.value.length > 0
                  ? [Math.floor(result.value[0]) * 1000, result.value[1]]
                  : [],
            } as LatestMetricRes);
          } else {
            return getErrorRes("No metric data found");
          }
        }
      });
      worker.on("exit", (code: any) => {
        removeKeyFromStorage(`${orgId}/${userId}/${type}/${from}`);
        logger.info(
          `WS_THREAD exited with code [${code}] for ${orgId}/${userId}/${type}`
        );
      });
    }

    return res;
  }

  @Subscription(() => LatestMetricRes, {
    topics: ({ args }) => `metric-${args.type}`,
    filter: ({ payload, args }) => {
      return args.nodeId === payload.nodeid;
    },
  })
  async getMetricRangeSub(
    @Root() payload: LatestMetricRes,
    @Args() args: SubMetricRangeInput
  ): Promise<LatestMetricRes> {
    await storeInStorage(
      `${args.orgId}/${args.userId}/${args.type}/${args.from}`,
      getTimestampCount("0")
    );
    return payload;
  }
}

export default MetricResolvers;
