/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Query, Resolver, Root, Subscription } from "type-graphql";
import { Worker } from "worker_threads";

import { METRIC_WS_INTERVAL } from "../../common/constants";
import {
  NotificationScopeEnumValue,
  NotificationTypeEnumValue,
  STATS_TYPE,
} from "../../common/enums";
import { logger } from "../../common/logger";
import { eventKeyToAction } from "../../common/notification";
import { addInStore, openStore } from "../../common/storage";
import {
  formatKPIValue,
  getBaseURL,
  getGraphsKeyByType,
  getScopesByRole,
  transformMetricsArray,
  wsUrlResolver,
} from "../../common/utils";
import {
  getNodeMetricRange,
  getNotifications,
} from "../datasource/subscriptions-api";
import { pubSub } from "./pubsub";
import {
  GetMetricBySiteInput,
  GetMetricByTabInput,
  GetMetricsSiteStatInput,
  GetMetricsStatInput,
  LatestMetricSubRes,
  MetricRes,
  MetricsRes,
  MetricsStateRes,
  NotificationsRes,
  NotificationsResDto,
  SubMetricsStatInput,
  SubSiteMetricByTabInput,
  SubSiteMetricsStatInput,
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

const processMetricResult = (metric: any) => ({
  msg: metric.msg,
  type: metric.type,
  success: metric.success,
  nodeId: metric.nodeId || "",
  siteId: metric.siteId || "",
  networkId: metric.networkId || "",
  packageId: metric.packageId || "",
  dataPlanId: metric.dataPlanId || "",
  value: metric.values[metric.values.length - 1][1],
});

const handleWebSocketMessage = (
  data: any,
  topic: string,
  pubSub: any,
  context?: { siteId?: string; nodeId?: string }
) => {
  if (!data.isError) {
    try {
      const res = JSON.parse(data.data);
      const results = Array.isArray(res.data.result)
        ? res.data.result
        : [res.data.result];

      results.forEach((result: any) => {
        if (result?.metric && result.value?.length === 2) {
          const metricSiteId = result?.metric?.siteid || context?.siteId || "";
          const metricNodeId = result?.metric?.nodeid || context?.nodeId || "";

          pubSub.publish(topic, {
            success: true,
            msg: "success",
            type: res.Name,
            nodeId: metricNodeId,
            siteId: metricSiteId,
            networkId: result?.metric?.network || "",
            packageId: result?.metric?.package || "",
            dataPlanId: result?.metric?.dataplan || "",
            value: [
              Math.floor(result.value[0]) * 1000,
              formatKPIValue(res.Name, result.value[1]),
            ],
          });
        }
      });
    } catch (error) {
      logger.error(`Failed to parse WebSocket message: ${error}`);
    }
  }
};

const setupWebSocketWorkers = (
  wsUrl: string,
  data: GetMetricsSiteStatInput,
  metricsKey: string[],
  siteIds: string[],
  nodeIds: string[],
  store: any
): Worker[] => {
  const workers: Worker[] = [];
  const sites = siteIds && siteIds.length > 0 ? siteIds : [""];
  const nodes = nodeIds && nodeIds.length > 0 ? nodeIds : [""];

  for (const siteId of sites) {
    for (const nodeId of nodes) {
      const topic = `stat-${data.orgName}-${data.userId}-${data.type}-${data.from}`;

      const workerContext = {
        siteId: siteId || "",
        nodeId: nodeId || "",
        topic: topic,
      };

      let url = `${wsUrl}/v1/live/metrics?interval=${METRIC_WS_INTERVAL}&operation=${
        data.operation || "mean"
      }&metric=${metricsKey.join(",")}`;

      if (nodeId) {
        url += `&node=${encodeURIComponent(nodeId)}`;
      }
      if (siteId) {
        url += `&site=${encodeURIComponent(siteId)}`;
      }

      const worker = new Worker(WS_THREAD, {
        workerData: {
          topic,
          url,
          context: workerContext,
        },
      });

      worker.on("message", (wsData: any) => {
        handleWebSocketMessage(wsData, topic, pubSub, workerContext);
      });

      worker.on("exit", async (code: any) => {
        await store.close();
        logger.info(`WS_THREAD exited with code [${code}] for ${topic}`);
      });

      workers.push(worker);
    }
  }

  return workers;
};
const fetchMetricsForSiteNodeCombination = async (
  baseURL: string,
  metricsKey: string[],
  siteId: string,
  nodeId: string,
  data: GetMetricsSiteStatInput
): Promise<any[]> => {
  const processedMetrics: any[] = [];
  const metricPromises = metricsKey.map(async key => {
    const combinedData = {
      ...data,
      siteIds: undefined,
      nodeIds: undefined,
      siteId: siteId || undefined,
      nodeId: nodeId || undefined,
    };

    try {
      const res = await getNodeMetricRange(baseURL, key, combinedData);
      if (Array.isArray(res.metrics)) {
        res.metrics.forEach(metric => {
          const processedMetric = processMetricResult(metric);
          if (!processedMetric.siteId && siteId) {
            processedMetric.siteId = siteId;
          }
          if (!processedMetric.nodeId && nodeId) {
            processedMetric.nodeId = nodeId;
          }
          processedMetrics.push(processedMetric);
        });
      }
    } catch (error) {
      logger.error(
        `Error fetching metrics for site ${siteId}, node ${nodeId}, key ${key}:`,
        error
      );
    }
  });

  await Promise.all(metricPromises);
  return processedMetrics;
};

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

    const wsUrl = wsUrlResolver(baseURL);
    const { type, from, userId, withSubscription, nodeId } = data;
    if (from === 0) throw new Error("Argument 'from' can't be zero.");

    const metricsKey = getGraphsKeyByType(type);
    const metrics: MetricsStateRes = { metrics: [] };

    if (metricsKey.length > 0) {
      const metricPromises = metricsKey.map(async key => {
        const res = await getNodeMetricRange(baseURL, key, { ...data });
        if (Array.isArray(res.metrics)) {
          res.metrics.forEach(metric => {
            metrics.metrics.push(processMetricResult(metric));
          });
        }
      });

      await Promise.all(metricPromises);
    }

    if (withSubscription && metrics.metrics.length > 0) {
      const topic = `${userId}-${type}-${from}`;
      const url = `${wsUrl}/v1/live/metrics?interval=${METRIC_WS_INTERVAL}&operation=${
        data.operation
      }&metric=${metricsKey.join(",")}${nodeId ? `&node=${nodeId}` : ""}${
        data.networkId ? `&network=${data.networkId}` : ""
      }`;

      const worker = new Worker(WS_THREAD, {
        workerData: { topic, url },
      });

      worker.on("message", (data: any) =>
        handleWebSocketMessage(data, topic, pubSub)
      );

      worker.on("exit", async (code: any) => {
        await store.close();
        logger.info(`WS_THREAD exited with code [${code}] for ${topic}`);
      });
    }
    metrics.metrics = metrics.metrics.filter(metric => metric.success === true);

    return metrics;
  }

  @Subscription(() => LatestMetricSubRes, {
    topics: ({ args }) => {
      return `${args.data.userId}-${args.data.type}-${args.data.from}`;
    },
  })
  async getMetricStatSub(
    @Root() payload: LatestMetricSubRes,
    @Arg("data") data: SubMetricsStatInput
  ): Promise<LatestMetricSubRes> {
    const store = openStore();
    await addInStore(store, `${data.userId}-${data.type}-${data.from}`, 0);
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
      logger.error(`Error getting base URL for metrics stat: ${baseURL}`);
      return { metrics: [] };
    }

    const { type, from, userId, nodeId, to } = data;
    if (from === 0) throw new Error("Argument 'from' can't be zero.");

    const metricsKey = getGraphsKeyByType(type);

    if (metricsKey.length > 0) {
      const metricPromises = metricsKey.map(
        async key =>
          await getNodeMetricRange(baseURL, key, {
            to,
            from,
            userId,
            nodeId,
            step: data.step,
            orgName: data.orgName,
            withSubscription: false,
            networkId: data.networkId,
            type: STATS_TYPE.ALL_NODE,
          })
      );

      const m = await Promise.all(metricPromises);
      return transformMetricsArray(m);
    }
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

    const wsUrl = wsUrlResolver(baseURL);

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

  @Query(() => MetricsStateRes)
  async getSiteStat(
    @Arg("data") data: GetMetricsSiteStatInput
  ): Promise<MetricsStateRes> {
    const store = openStore();

    try {
      const { message: baseURL, status } = await getBaseURL(
        "metrics",
        data.orgName,
        store
      );
      if (status !== 200) {
        logger.error(`Error getting base URL for metrics stat: ${baseURL}`);
        return { metrics: [] };
      }

      const wsUrl = wsUrlResolver(baseURL);
      const { type, from, withSubscription, nodeIds, siteIds } = data;
      if (from === 0) throw new Error("Argument 'from' can't be zero.");

      const metricsKey = getGraphsKeyByType(type);
      const metrics: MetricsStateRes = { metrics: [] };

      if (metricsKey.length > 0) {
        const sites = siteIds && siteIds.length > 0 ? siteIds : [""];
        const nodes = nodeIds && nodeIds.length > 0 ? nodeIds : [""];

        for (const siteId of sites) {
          for (const nodeId of nodes) {
            const processedMetrics = await fetchMetricsForSiteNodeCombination(
              baseURL,
              metricsKey,
              siteId,
              nodeId,
              data
            );
            metrics.metrics.push(...processedMetrics);
          }
        }
      }

      if (withSubscription && metrics.metrics.length > 0) {
        setupWebSocketWorkers(
          wsUrl,
          data,
          metricsKey,
          siteIds || [],
          nodeIds || [],
          store
        );
      } else {
        await store.close();
      }

      return metrics;
    } catch (error) {
      await store.close();
      throw error;
    }
  }

  @Subscription(() => LatestMetricSubRes, {
    topics: ({ args }) => {
      return `stat-${args.data.orgName}-${args.data.userId}-${args.data.type}-${args.data.from}`;
    },
  })
  async getSiteMetricStatSub(
    @Root() payload: LatestMetricSubRes,
    @Arg("data") data: SubSiteMetricsStatInput
  ): Promise<LatestMetricSubRes> {
    const store = openStore();
    await addInStore(
      store,
      `stat-${data.orgName}-${data.userId}-${data.type}-${data.from}`,
      0
    );
    await store.close();
    return payload;
  }

  @Query(() => MetricsRes)
  async getMetricBySite(@Arg("data") data: GetMetricBySiteInput) {
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

    const { type, from, userId, siteId, to } = data;
    if (from === 0) throw new Error("Argument 'from' can't be zero.");

    const metricsKey = getGraphsKeyByType(type);

    if (metricsKey.length > 0) {
      const metricPromises = metricsKey.map(
        async key =>
          await getNodeMetricRange(baseURL, key, {
            to,
            from,
            userId,
            siteId,
            step: data.step,
            orgName: data.orgName,
            withSubscription: false,
            type: STATS_TYPE.SITE,
          })
      );

      const m = await Promise.all(metricPromises);
      return transformMetricsArray(m);
    }
  }

  @Subscription(() => LatestMetricSubRes, {
    topics: ({ args }) => {
      return `${args.data.userId}/${args.data.type}/${args.data.from}`;
    },
  })
  async getSiteMetricByTabSub(
    @Root() payload: LatestMetricSubRes,
    @Arg("data") data: SubSiteMetricByTabInput
  ): Promise<LatestMetricSubRes> {
    const store = openStore();
    await addInStore(store, `${data.userId}/${payload.type}/${data.from}`, 0);
    await store.close();
    return payload;
  }
}

export default SubscriptionsResolvers;
