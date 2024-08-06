import { Arg, Args, Query, Resolver, Root, Subscription } from "type-graphql";
import { Worker } from "worker_threads";

import {
  METRIC_API_GW_SOCKET,
  NOTIFICATION_API_GW_WS,
  STORAGE_KEY,
} from "../../common/configs";
import { logger } from "../../common/logger";
import { addInStore, openStore, removeFromStore } from "../../common/storage";
import { getGraphsKeyByType, getTimestampCount } from "../../common/utils";
import {
  getNodeRangeMetric,
  getNotifications,
} from "../datasource/subscriptions-api";
import { pubSub } from "./pubsub";
import {
  GetMetricByTabInput,
  GetNotificationsInput,
  LatestMetricRes,
  MetricRes,
  MetricsRes,
  NotificationsRes,
  NotificationsResDto,
  SubMetricByTabInput,
} from "./types";

const WS_THREAD = "./threads/MetricsWSThread.js";
const NOTIFICATION_THREAD = "./threads/NotificationsWSThread.js";

const getErrorRes = (msg: string) =>
  ({
    orgId: "",
    msg: msg,
    type: "",
    nodeId: "",
    values: [],
    success: false,
  } as MetricRes);

@Resolver(String)
class SubscriptionsResolvers {
  @Query(() => MetricsRes)
  async getMetricByTab(@Arg("data") data: GetMetricByTabInput) {
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
          removeFromStore(openStore(), `${orgId}/${userId}/${type}/${from}`);
          logger.info(
            `WS_THREAD exited with code [${code}] for ${orgId}/${userId}/${type}`
          );
        });
      });
    }
    return metrics;
  }

  @Query(() => NotificationsRes)
  async getNotifications(@Arg("data") data: GetNotificationsInput) {
    const notifications = getNotifications(data);
    const workerData = {
      url: `${NOTIFICATION_API_GW_WS}/v1/distributor/live`,
      orgId: data.orgId,
      scope: data.scopes,
      userId: data.userId,
      networkId: data.networkId,
      key: "UKAMA_NOTIFICATION_STORAGE_KEY",
    };

    const worker = new Worker(NOTIFICATION_THREAD, {
      workerData,
    });

    const key = `notification-${data.userId}-${data.orgId}-${data.forRole}-${data.networkId}`;
    worker.on("message", (_data: any) => {
      if (!_data.isError) {
        const res = JSON.parse(_data.data);
        if (res && res.id) {
          pubSub.publish(key, {
            createdAt: res.createdAt,
            description: res.description,
            id: res.id,
            isRead: res.isRead,
            scope: res.scope,
            title: res.title,
            type: res.type,
          } as NotificationsResDto);
        } else {
          return getErrorRes("No notification data found");
        }
      }
    });

    worker.on("exit", (code: any) => {
      removeFromStore(openStore(), key);
      logger.info(
        `WS_THREAD exited with code [${code}] for ${data.orgId}/${data.userId}/${data.networkId}`
      );
    });
    return notifications;
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
    await addInStore(
      openStore(),
      `${args.orgId}/${args.userId}/${payload.type}/${args.from}`,
      getTimestampCount("0")
    );
    return payload;
  }

  @Subscription(() => NotificationsResDto, {
    topics: ({ args }) => {
      return `notification-${args.userId}-${args.orgId}-${args.forRole}-${args.networkId}`;
    },
  })
  async notificationSubscription(
    @Root() payload: NotificationsResDto,
    @Args() args: GetNotificationsInput
  ): Promise<NotificationsResDto> {
    await addInStore(
      openStore(),
      `notification-${args.userId}-${args.orgId}-${args.forRole}-${args.networkId}`,
      getTimestampCount("0")
    );
    return payload;
  }
}

export default SubscriptionsResolvers;
