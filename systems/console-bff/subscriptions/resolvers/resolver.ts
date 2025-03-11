import { Arg, Query, Resolver, Root, Subscription } from "type-graphql";
import { Worker } from "worker_threads";

import {
  NotificationScopeEnumValue,
  NotificationTypeEnumValue,
} from "../../common/enums";
import { logger } from "../../common/logger";
import { addInStore, openStore } from "../../common/storage";
import {
  eventKeyToAction,
  getBaseURL,
  getGraphsKeyByType,
  getScopesByRole,
} from "../../common/utils";
import {
  getMetricRange,
  getNodeRangeMetric,
  getNotifications,
} from "../datasource/subscriptions-api";
import { pubSub } from "./pubsub";
import {
  GetMetricByTabInput,
  GetMetricsStatInput,
  LatestMetricSubRes,
  MetricRes,
  MetricStateRes,
  MetricsRes,
  MetricsStateRes,
  NotificationsRes,
  NotificationsResDto,
  SubMetricByTabInput,
  SubMetricsStatInput,
} from "./types";

const WS_THREAD = "./threads/MetricsWSThread.mjs";
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
  @Query(() => MetricsStateRes)
  async getMetricsStat(
    @Arg("data") data: GetMetricsStatInput
  ): Promise<MetricsStateRes> {
    const store = openStore();
    const { message: baseURL, status } = await getBaseURL(
      "metrics",
      data.orgName,
      store
    );
    if (status !== 200) {
      logger.error(`Error getting base URL for metrics stat: ${baseURL}`);
      return { metrics: [] };
    }

    let wsUrl = baseURL;
    if (wsUrl?.includes("https://")) {
      wsUrl = wsUrl.replace("https://", "wss://");
    } else if (wsUrl?.startsWith("http://")) {
      wsUrl = wsUrl.replace("http://", "ws://");
    }

    const { type, from, userId, withSubscription } = data;
    if (from === 0) throw new Error("Argument 'from' can't be zero.");

    const metricsKey = getGraphsKeyByType(type);
    const metrics: MetricsStateRes = { metrics: [] };

    if (metricsKey.length > 0) {
      const metricPromises = metricsKey.map(async key => {
        const res = await getMetricRange(baseURL, key, { ...data });
        let avg = 0;

        for (let i = 0; i < res.values.length; i++) {
          if (res.values[i][1] === 0) {
            res.values.splice(i, 1);
          }
        }

        if (Array.isArray(res.values) && res.values.length > 0) {
          if (
            res.values.length === 1 ||
            key === "unit_uptime" ||
            key === "network_uptime"
          ) {
            avg = res.values[res.values.length - 1][1];
          } else {
            const sum = res.values.reduce((acc, val) => acc + val[1], 0);
            avg = sum / res.values.length;
          }
        }

        return {
          msg: res.msg,
          type: res.type,
          nodeId: res.nodeId,
          success: res.success,
          value: parseFloat(Number(avg).toFixed(2)),
        };
      });

      metrics.metrics = await Promise.all(metricPromises);
    }

    if (withSubscription && metrics.metrics.length > 0) {
      metrics.metrics.forEach((metric: MetricStateRes) => {
        const workerData = {
          topic: `stat-sub-${userId}/${type}/${from}`,
          url: `${wsUrl}/v1/live/metrics?interval=${data.step}&metric=${metric.type}&node=${metric.nodeId}`,
        };

        const worker = new Worker(WS_THREAD, {
          workerData,
        });

        worker.on("message", (_data: any) => {
          if (!_data.isError) {
            const res = JSON.parse(_data.data);
            const result = res.data.result[0];
            if (result && result.metric && result.value.length > 0) {
              pubSub.publish(workerData.topic, {
                success: true,
                msg: "success",
                type: metric.type,
                nodeId: metric.nodeId,
                value: [
                  Math.floor(result.value[0]) * 1000,
                  parseFloat(Number(result.value[1] || 0).toFixed(2)),
                ],
              });
            }
          }
        });
        worker.on("exit", async (code: any) => {
          await store.close();
          logger.info(
            `WS_THREAD exited with code [${code}] for ${userId}/${type}/${from}`
          );
        });
      });
    }

    return metrics;
  }

  @Subscription(() => LatestMetricSubRes, {
    topics: ({ args }) => {
      return `stat-sub-${args.data.userId}/${args.data.type}/${args.data.from}`;
    },
  })
  async getMetricStatSub(
    @Root() payload: LatestMetricSubRes,
    @Arg("data") data: SubMetricsStatInput
  ): Promise<LatestMetricSubRes> {
    const store = openStore();
    await addInStore(
      store,
      `stat-sub-${data.userId}/${payload.type}/${data.from}`,
      0
    );
    await store.close();
    return payload;
  }

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

    let wsUrl = baseURL;
    if (wsUrl?.includes("https://")) {
      wsUrl = wsUrl.replace("https://", "wss://");
    } else if (wsUrl?.startsWith("http://")) {
      wsUrl = wsUrl.replace("http://", "ws://");
    }

    const { type, from, nodeId, userId, withSubscription } = data;
    if (from === 0) throw new Error("Argument 'from' can't be zero.");
    const metricsKey: string[] = getGraphsKeyByType(type);
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
      subKey.split(",").forEach((key: string) => {
        if (key === "") return;
        const workerData = {
          topic: `${userId}/${key}/${from}`,
          url: `${wsUrl}/v1/live/metrics?interval=${data.step}&metric=${key}&node=${nodeId}`,
        };
        const worker = new Worker(WS_THREAD, {
          workerData,
        });

        worker.on("message", (_data: any) => {
          if (!_data.isError) {
            const res = JSON.parse(_data.data);
            const result = res.data.result[0];
            if (result && result.metric && result.value.length > 0) {
              pubSub.publish(workerData.topic, {
                type: key,
                success: true,
                msg: "success",
                nodeId: nodeId,
                value: [
                  Math.floor(result.value[0]) * 1000,
                  parseFloat(Number(result.value[1] || 0).toFixed(2)),
                ],
              });
            }
          }
        });
        worker.on("exit", async (code: any) => {
          await store.close();
          logger.info(
            `WS_THREAD exited with code [${code}] for ${userId}/${type}/${from}`
          );
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
      await store.close();
      logger.info(
        `WS_THREAD exited with code [${code}] for ${orgId}/${userId}/${networkId}/${subscriberId}/${startTimestamp}`
      );
      worker.terminate();
    });
    return notifications;
  }

  @Subscription(() => LatestMetricSubRes, {
    topics: ({ args }) => {
      return `${args.data.userId}/${args.data.type}/${args.data.from}`;
    },
  })
  async getMetricByTabSub(
    @Root() payload: LatestMetricSubRes,
    @Arg("data") data: SubMetricByTabInput
  ): Promise<LatestMetricSubRes> {
    const store = openStore();
    await addInStore(store, `${data.userId}/${payload.type}/${data.from}`, 0);
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
