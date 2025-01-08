import { Arg, Query, Resolver, Root, Subscription } from "type-graphql";
import { Worker } from "worker_threads";

import { STORAGE_KEY } from "../../common/configs";
import {
  NotificationScopeEnumValue,
  NotificationTypeEnumValue,
} from "../../common/enums";
import { logger } from "../../common/logger";
import { addInStore, openStore, removeFromStore } from "../../common/storage";
import {
  eventKeyToAction,
  getBaseURL,
  getGraphsKeyByType,
  getScopesByRole,
  getTimestampCount,
} from "../../common/utils";
import {
  getNodeRangeMetric,
  getNotifications,
} from "../datasource/subscriptions-api";
import { pubSub } from "./pubsub";
import {
  GetMetricByTabInput,
  LatestMetricRes,
  MetricRes,
  MetricsRes,
  NotificationsRes,
  NotificationsResDto,
  SubMetricByTabInput,
} from "./types";

const WS_THREAD = "./threads/MetricsWSThread.js";
const NOTIFICATION_THREAD = "./threads/NotificationsWSThread.mjs";

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
    const store = openStore();
    const { message: baseURL, status } = await getBaseURL(
      "metrics",
      data.orgName,
      store
    );
    if (status !== 200) {
      logger.error(`Error getting base URL for notification: ${baseURL}`);
      return { notifications: [] };
    }
    const { type, orgId, userId, nodeId, withSubscription, from } = data;
    if (from === 0) throw new Error("Argument 'from' can't be zero.");
    const metricsKey: string[] = getGraphsKeyByType(type, nodeId);
    const metrics: MetricsRes = { metrics: [] };
    if (metricsKey.length > 0) {
      for (let i = 0; i < metricsKey.length; i++) {
        const res = await getNodeRangeMetric(baseURL, {
          ...data,
          type: metricsKey[i],
        });
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
          url: `${baseURL}/v1/live/metrics?interval=1&metric=${key}&node=${nodeId}`,
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
        worker.on("exit", async (code: any) => {
          await removeFromStore(store, `${orgId}/${userId}/${type}/${from}`);
          logger.info(
            `WS_THREAD exited with code [${code}] for ${orgId}/${userId}/${type}`
          );
          await store.close();
        });
      });
    }
    return metrics;
  }

  @Query(() => NotificationsRes)
  async getNotifications(
    @Arg("orgId") orgId: string,
    @Arg("role") role: string,
    @Arg("userId") userId: string,
    @Arg("orgName") orgName: string,
    @Arg("networkId") networkId: string,
    @Arg("subscriberId") subscriberId: string,
    @Arg("startTimestamp") startTimestamp: string
  ) {
    const store = openStore();
    const { message: baseURL, status } = await getBaseURL(
      "notification",
      orgName,
      store
    );
    if (status !== 200) {
      logger.error(`Error getting base URL for notification: ${baseURL}`);
      return { notifications: [] };
    }

    let wsUrl = baseURL;
    if (wsUrl?.includes("https://")) {
      wsUrl = wsUrl.replace("https://", "wss://");
    } else if (wsUrl?.startsWith("http://")) {
      wsUrl = wsUrl.replace("http://", "ws://");
    }

    const notifications = getNotifications(
      baseURL,
      orgId,
      userId,
      networkId,
      subscriberId
    );
    const scopesPerRole = getScopesByRole(role);
    let scopes = "";
    if (scopesPerRole.length > 0) {
      for (const scope of scopesPerRole) {
        scopes = scopes + `&scope=${scope}`;
      }
      scopes = scopes.substring(1);
    }

    const key = `notification-${orgId}-${userId}-${networkId}-${subscriberId}-${startTimestamp}`;
    const workerData = {
      url: `${wsUrl}/v1/distributor/live`,
      key: key,
      orgId: orgId,
      scopes: scopes,
      userId: userId,
      networkId: networkId,
      subscriberId: subscriberId,
    };

    const worker = new Worker(NOTIFICATION_THREAD, {
      workerData,
    });

    worker.on("message", (_data: any) => {
      if (!_data.isError) {
        const res = JSON.parse(_data.data);
        if (res && res.id) {
          const n: NotificationsResDto = {
            id: res.id,
            isRead: false,
            title: res.title,
            eventKey: res.eventKey,
            createdAt: res.createdAt,
            resourceId: res.resourceId,
            description: res.description,
            type: NotificationTypeEnumValue(res.type),
            scope: NotificationScopeEnumValue(res.scope),
          };
          n.redirect = eventKeyToAction(res.event_key, n);
          pubSub.publish(key, n);
        } else {
          return getErrorRes("No notification data found");
        }
      }
    });

    worker.on("exit", async (code: any) => {
      await removeFromStore(store, key);
      logger.info(
        `WS_THREAD exited with code [${code}] for ${orgId}/${userId}/${networkId}/${subscriberId}/${startTimestamp}`
      );
      store.close();
      worker.terminate();
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
    @Arg("data") data: SubMetricByTabInput
  ): Promise<LatestMetricRes> {
    const store = openStore();
    await addInStore(
      store,
      `${data.orgId}/${data.userId}/${payload.type}/${data.from}`,
      getTimestampCount("0")
    );
    await store.close();
    return payload;
  }

  @Subscription(() => NotificationsResDto, {
    topics: ({ args }) => {
      return `notification-${args.orgId}-${args.userId}-${args.networkId}-${args.subscriberId}-${args.startTimestamp}`;
    },
  })
  async notificationSubscription(
    @Root() payload: NotificationsResDto,
    @Arg("orgId") orgId: string,
    @Arg("role") role: string,
    @Arg("userId") userId: string,
    @Arg("orgName") orgName: string,
    @Arg("networkId") networkId: string,
    @Arg("subscriberId") subscriberId: string,
    @Arg("startTimestamp") startTimestamp: string
  ): Promise<NotificationsResDto> {
    const store = openStore();
    await addInStore(
      store,
      `notification-${orgId}-${userId}-${networkId}-${subscriberId}-${startTimestamp}`,
      0
    );
    logger.info("Notification payload :", payload);
    await store.close();
    return payload;
  }
}

export default SubscriptionsResolvers;
