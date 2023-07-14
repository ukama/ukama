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

import { METRIC_API_GW, STORAGE_KEY } from "../../../common/configs";
import { logger } from "../../../common/logger";
import { storeInStorage } from "../../../common/storage";
import { getTimestampCount } from "../../../common/utils";
import { GetMetricInput, MetricRes } from "./types";

const WS_THREAD = "./threads/MetricsWSThread.ts";

@Resolver(MetricRes)
class MetricResolvers {
  @Query(() => MetricRes)
  async getMetrics(
    @Arg("data") data: GetMetricInput,
    @PubSub() pubSub: PubSubEngine
  ) {
    const { type, orgId, userId, nodeId } = data;
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
    worker.on("message", (_data: any) => {
      if (!_data.isError) {
        const res = JSON.parse(_data.data);
        pubSub.publish(`metric-${type}`, {
          env: res.metric.env,
          nodeid: res.metric.nodeid,
          type: res.metric.type,
          value: [{ x: res.value[0][0], y: res.value[0][1] }],
        } as MetricRes);
      }
    });
    worker.on("exit", (code: any) => {
      logger.info(
        `WS_THREAD exited with code [${code}] for ${orgId}/${userId}/${type}`
      );
    });
    const m: MetricRes = {
      env: "dev",
      nodeid: "uk-test17-hnode-a1-31df",
      type: "node",
      value: [{ x: 1687397186.619, y: 2165.30859375 }],
    };
    return m;
  }

  @Subscription(() => MetricRes, {
    topics: ({ args }) => `metric-${args.type}`,
    filter: ({ payload, args }) => {
      return args.nodeId === payload.nodeid;
    },
  })
  async getMetricEvent(
    @Root() payload: MetricRes,
    @Args() args: GetMetricInput
  ): Promise<MetricRes> {
    await storeInStorage(
      `${args.orgId}/${args.userId}/${args.type}`,
      getTimestampCount("0")
    );
    return payload;
  }
}

export default MetricResolvers;
