import WebSocket from "ws";

import { METRIC_API_GW } from "../../../common/configs";
import { Resolvers } from "../../types";
import {
  removeKeyFromBucket,
  retriveFromBucket,
  storeInBucket,
} from "../bucket";

const getTimestamp = (count: string) =>
  parseInt((Date.now() / 1000).toString()) + "-" + count;

const metricResolvers: Resolvers = {
  Query: {
    getMetrics: async (_, { input }, { pubSub }) => {
      const { type, orgId, userId } = input;
      const ws = new WebSocket(METRIC_API_GW);
      ws.on("error", e => console.log(e));

      ws.on("open", async function open() {
        await storeInBucket(`${orgId}/${userId}/${type}`, getTimestamp("0"));
      });

      ws.on("message", async function message(data) {
        const value = await retriveFromBucket(`${orgId}/${userId}/${type}`);
        let occurance = value ? parseInt(value.split("-")[1]) : 0;
        occurance += 1;

        await storeInBucket(
          `${orgId}/${userId}/${type}`,
          getTimestamp(`${occurance}`)
        );

        if (occurance === 9) {
          ws.close();
          ws.terminate();
          removeKeyFromBucket(`${orgId}/${userId}/${type}`);
        }

        pubSub.publish(`metric-${input.type}`, `${orgId}/${userId}/${type}`, {
          value: data.toString(),
        });
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
        await storeInBucket(
          `${input.orgId}/${input.userId}/${input.type}`,
          getTimestamp("0")
        );
        return payload;
      },
    },
  },
};

export default metricResolvers;
