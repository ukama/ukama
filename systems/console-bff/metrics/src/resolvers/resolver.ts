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
          value: res.value,
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
      value: [
        [1690211392, 1.935],
        [1690211392, 1.935],
        [1690211392, 1.935],
        [1690211392, 0.808],
        [1690211392, 0.808],
        [1690211393, 0.729],
        [1690211393, 0.729],
        [1690211393, 1.36],
        [1690211394, 0.597],
        [1690211394, 1.838],
        [1690211395, 1.703],
        [1690211395, 0.32],
        [1690211396, 1.926],
        [1690211396, 0.49],
        [1690211397, 1.115],
        [1690211397, 0.798],
        [1690211398, 1.276],
        [1690211398, 1.975],
        [1690211399, 1.657],
        [1690211399, 1.044],
        [1690211400, -0.24],
        [1690211400, 1.36],
        [1690211401, 0.038],
        [1690211401, 1.662],
        [1690211402, 1.161],
        [1690211402, -0.2],
        [1690211403, -0.1],
        [1690211403, 0.165],
        [1690211404, 0.251],
        [1690211404, 0.619],
        [1690211405, -0.05],
        [1690211405, 2.194],
        [1690211406, 0.433],
        [1690211406, 1.57],
        [1690211407, 0.56],
        [1690211407, 1.968],
        [1690211408, 0.773],
        [1690211408, 1.363],
        [1690211409, 1.333],
        [1690211409, 1.519],
        [1690211410, 0.492],
        [1690211410, 2.048],
        [1690211411, 1.013],
        [1690211411, 2.142],
        [1690211412, 0.427],
        [1690211412, 0.702],
        [1690211413, -0.04],
      ],
    };
    return m;
  }

  @Subscription(() => MetricRes, {
    topics: ({ args }) => `metric-${args.type}`,
    filter: ({ payload, args }) => {
      return args.nodeId === payload.nodeid;
    },
  })
  async getMetric(
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
