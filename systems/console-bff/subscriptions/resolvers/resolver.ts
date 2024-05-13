/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import crypto from "crypto";
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
import { getGraphsKeyByType, getTimestampCount } from "../../common/utils";
import { getLatestMetric, getNodeRangeMetric } from "../datasource/metrics-api";
import {
  GetLatestMetricInput,
  GetMetricByTabInput,
  LatestMetricRes,
  MetricRes,
  MetricsRes,
  StatsMetric,
  SubMetricByTabInput,
} from "./types";

const WS_THREAD = "./threads/MetricsWSThread.js";

const getErrorRes = (msg: string) =>
  ({
    orgId: "",
    msg: msg,
    type: "",
    nodeId: "",
    values: [],
    success: false,
  } as MetricRes);

/*
interface WorkerData {
  type: string;
  orgId: string;
  userId: string;
  timestamp: number;
  key: string;
  url: string;
}

const constructUrl = (type: string): string => {
  return `${METRIC_API_GW_SOCKET}/v1/live/metrics?interval=1&metric=${type}`;
};

const createWorkerData = (
  type: string,
  orgId: string,
  userId: string,
  from: number
): WorkerData => {
  return {
    type,
    orgId,
    userId,
    timestamp: from,
    key: STORAGE_KEY,
    url: constructUrl(type),
  };
};

const handleWorkerMessage = (
  _data: any,
  type: string,
  nodeId: string,
  pubSub: any
) => {
  if (!_data.isError) {
    const res = JSON.parse(_data.data);
    const result = res.data.result[0];
    if (
      result &&
      result.metric &&
      Array.isArray(result.value) &&
      result.value.length > 0
    ) {
      pubSub.publish(`metric-${type}`, {
        orgId: result.metric.org,
        nodeId: nodeId,
        type: type,
        value: [Math.floor(result.value[0]) * 1000, result.value[1]],
      } as LatestMetricRes);
    } else {
      return getErrorRes("No metric data found");
    }
  }
};

@Query(() => MetricRes)
  async getMetricRange(
    @Arg("data") data: GetMetricRangeInput,
    @PubSub() pubSub: PubSubEngine
  ) {
    const { type, orgId, userId, nodeId, withSubscription, from } = data;
    if (from === 0) throw new Error("Argument 'from' can't be zero.");
    const res = await getMetricRange(data);
    if (withSubscription && res.orgId && res.nodeId) {
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
              orgId: result.metric.org,
              nodeId: nodeId,
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
    const {
      type,
      orgId = "",
      userId = "",
      nodeId,
      withSubscription,
      from,
    } = data;
    if (from === 0) throw new Error("Argument 'from' can't be zero.");
    const res = await getNodeRangeMetric(data);

    if (withSubscription && res.orgId && res.nodeId) {
      const workerData = createWorkerData(type, orgId, userId, from);
      const worker = new Worker(WS_THREAD, {
        workerData,
      });
      worker.on("message", (_data: any) =>
        handleWorkerMessage(_data, type, nodeId, pubSub)
      );
      worker.on("error", err => {
        logger.error(`Worker error: ${err}`);
      });
      worker.on("exit", (code: any) => {
        const keys = getGraphsKeyByType(type, nodeId);
        keys.forEach(async (key: string) => {
          await removeKeyFromStorage(`${orgId}/${userId}/${key}/${from}`);
        });
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
    return args.nodeId === payload.nodeId;
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
*/

@Resolver(MetricRes)
class MetricResolvers {
  @Query(() => StatsMetric)
  async getStatsMetric() {
    return {
      activeSubscriber: Math.floor(crypto.randomInt(1, 30)),
      averageThroughput: Math.floor(crypto.randomInt(1, 50)),
      averageSignalStrength: Math.floor(crypto.randomInt(1, 90)),
    };
  }

  @Query(() => LatestMetricRes)
  async getLatestMetric(@Arg("data") data: GetLatestMetricInput) {
    return await getLatestMetric(data);
  }

  @Query(() => MetricsRes)
  async getMetricByTab(
    @Arg("data") data: GetMetricByTabInput,
    @PubSub() pubSub: PubSubEngine
  ) {
    const { type, orgId, userId, nodeId, withSubscription, from } = data;
    if (from === 0) throw new Error("Argument 'from' can't be zero.");
    const metricsKey: string[] = getGraphsKeyByType(type, nodeId);
    const metrics: MetricsRes = { metrics: [] };
    if (metricsKey.length > 0) {
      for (let i = 0; i < metricsKey.length; i++) {
        const res = await getNodeRangeMetric({ ...data, type: metricsKey[i] });
        metrics.metrics.push(res);
      }
    }
    if (withSubscription && metrics.metrics.length > 0) {
      let subKey = "";
      metrics.metrics.forEach((metric: MetricRes) => {
        if (metric.values.length > 2) subKey = subKey + metric.type + ",";
      });
      subKey = subKey.slice(0, -1);
      subKey.split(",").forEach((key: string) => {
        const workerData = {
          type: key,
          orgId,
          userId,
          url: `${METRIC_API_GW_SOCKET}/v1/live/metrics?interval=1&metric=${key}&node=${nodeId}`,
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
              pubSub.publish(key, {
                success: true,
                msg: "success",
                orgId: result.metric.org,
                nodeId: nodeId,
                type: key,
                userId: userId,
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
      });
    }
    return metrics;
  }

  @Subscription(() => LatestMetricRes, {
    topics: ({ args }) => {
      return getGraphsKeyByType(args.type, args.nodeId);
    },
    filter: ({ payload, args }) => {
      return args.nodeId === payload.nodeId && args.userId === payload.userId;
    },
  })
  async getMetricByTabSub(
    @Root() payload: LatestMetricRes,
    @Args() args: SubMetricByTabInput
  ): Promise<LatestMetricRes> {
    await storeInStorage(
      `${args.orgId}/${args.userId}/${payload.type}/${args.from}`,
      getTimestampCount("0")
    );
    return payload;
  }
}

export default MetricResolvers;
