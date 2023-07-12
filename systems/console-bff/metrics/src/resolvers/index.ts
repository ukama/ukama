import { Worker } from "worker_threads";

import { METRIC_API_GW, STORAGE_KEY } from "../../../common/configs";
import { logger } from "../../../common/logger";
import { storeInStorage } from "../../../common/storage";
import { getTimestampCount } from "../../../common/utils/";
import { Metric, Resolvers } from "../../types";

const WS_THREAD = "./threads/MetricsWSThread.ts";

const metricResolvers: Resolvers = {
  Query: {
    getMetrics: async (_, { input }, { pubSub }) => {
      const { type, orgId, userId } = input;
      const workerData: any = {
        type,
        orgId,
        userId,
        url: METRIC_API_GW,
        key: STORAGE_KEY,
      };
      const worker = new Worker(WS_THREAD, {
        workerData,
      });
      worker.on("message", (data: any) => {
        if (!data.isError) {
          const res = JSON.parse(data.data);
          pubSub.publish(`metric-${input.type}`, `${orgId}/${userId}/${type}`, {
            env: res.metric.env,
            nodeid: res.metric.nodeid,
            type: res.metric.type,
            value: res.value,
          } as Metric);
        }
      });
      worker.on("exit", (code: any) => {
        logger.info(
          `WS_THREAD exited with code [${code}] for ${orgId}/${userId}/${type}`
        );
      });
      const m: Metric = {
        env: "dev",
        nodeid: "uk-test17-hnode-a1-31df",
        type: "node",
        value: [[1687397186.619, 2165.30859375]],
      };
      return m;
    },
  },
  Subscription: {
    getMetricEvent: {
      subscribe: async (_, { input: { orgId, userId, type } }, { pubSub }) =>
        pubSub.subscribe(`metric-${type}`, `${orgId}/${userId}/${type}`),
      resolve: async (payload, { input }): Promise<Metric> => {
        await storeInStorage(
          `${input.orgId}/${input.userId}/${input.type}`,
          getTimestampCount("0")
        );

        return payload;
      },
    },
  },
};

export default metricResolvers;
