import { Worker } from "worker_threads";

import { METRIC_API_GW, STORAGE_KEY } from "../../../common/configs";
import { logger } from "../../../common/logger";
import { storeInStorage } from "../../../common/storage";
import { getTimestampCount } from "../../../common/utils/";
import { Resolvers } from "../../types";

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
          pubSub.publish(`metric-${input.type}`, `${orgId}/${userId}/${type}`, {
            value: data.data,
          });
        }
      });
      worker.on("exit", (code: any) => {
        logger.info(
          `WS_THREAD exited with code [${code}] for ${orgId}/${userId}/${type}`
        );
      });
      return [
        {
          value: "1",
        },
      ];
    },
  },
  Subscription: {
    getMetricEvent: {
      subscribe: async (_, { input: { orgId, userId, type } }, { pubSub }) =>
        pubSub.subscribe(`metric-${type}`, `${orgId}/${userId}/${type}`),
      resolve: async (payload, { input }) => {
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
