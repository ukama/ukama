/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
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
  formatKPIValue,
  getBaseURL,
  getGraphsKeyByType,
  getScopesByRole,
  wsUrlResolver,
} from "../../common/utils";
import {
  getNodeMetricRange,
  getNotifications,
  getSiteMetricRange,
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
  SiteMetricsStateRes,
  SubMetricByTabInput,
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
        let avg = 0;

        res.values = res.values.filter(value => value[1] !== 0);

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
          nodeId: res.nodeId || "",
          success: res.success,
          value: formatKPIValue(res.type, avg),
        };
      });

      metrics.metrics = await Promise.all(metricPromises);
    }

    if (withSubscription && metrics.metrics.length > 0) {
      const workerData = {
        topic: `${userId}-${type}-${from}`,
        url: `${wsUrl}/v1/live/metrics?interval=${
          data.step
        }&metric=${metricsKey.join(",")}&node=${nodeId}`,
      };
      const worker = new Worker(WS_THREAD, {
        workerData,
      });

      worker.on("message", (_data: any) => {
        if (!_data.isError) {
          try {
            const res = JSON.parse(_data.data);
            const result = res.data.result[0];
            if (result && result.metric && result.value.length > 0) {
              pubSub.publish(workerData.topic, {
                success: true,
                msg: "success",
                type: res.Name,
                nodeId: data.nodeId,
                value: [
                  Math.floor(result.value[0]) * 1000,
                  formatKPIValue(res.Name, result.value[1]),
                ],
              });
            }
          } catch (error) {
            logger.error(`Failed to parse WebSocket message: ${error}`);
          }
        }
      });
      worker.on("exit", async (code: any) => {
        await store.close();
        logger.info(
          `WS_THREAD exited with code [${code}] for ${userId}/${type}/${from}`
        );
      });
    }
    metrics.metrics = metrics.metrics.filter(metric => metric.success === true);

    return metrics;
  }

  @Query(() => SiteMetricsStateRes)
  async getSiteStat(
    @Arg("data") data: GetMetricsSiteStatInput
  ): Promise<SiteMetricsStateRes> {
    const store = openStore();
    const { message: baseURL, status } = await getBaseURL(
      "metrics",
      data.orgName,
      store
    );
    if (status !== 200) {
      logger.error(`Error getting base URL for site stat: ${baseURL}`);
      return { metrics: [] };
    }

    const wsUrl = wsUrlResolver(baseURL);
    const { from, userId, withSubscription, siteIds, type, nodeIds } = data;
    if (from === 0) throw new Error("Argument 'from' can't be zero.");
    if (!siteIds || siteIds.length === 0) {
      throw new Error("At least one siteId must be provided");
    }

    const metrics: SiteMetricsStateRes = { metrics: [] };
    const metricKeys = getGraphsKeyByType(type);

    for (const siteId of siteIds) {
      try {
        const combinedResult = await Promise.all(
          metricKeys.map(async key => {
            try {
              return await getSiteMetricRange(baseURL, key, {
                ...data,
                siteId,
              });
            } catch (error) {
              logger.error(
                `Error processing site metric ${key} for site ${siteId}: ${error}`
              );
              return {
                msg: "Error",
                type: key,
                siteId: siteId || "",
                nodeId: "",
                success: false,
                values: [],
              };
            }
          })
        );

        for (const res of combinedResult) {
          let avg = 0;

          if (Array.isArray(res.values)) {
            res.values = res.values.filter(value => value[1] !== 0);
            if (res.values.length === 1 || res.type === "site_uptime_seconds") {
              avg = res.values[res.values.length - 1][1];
            } else if (res.type === "site_uptime_percentage") {
              const sum = res.values.reduce((acc, val) => acc + val[1], 0);
              avg = sum / res.values.length;
            } else {
              const sum = res.values.reduce((acc, val) => acc + val[1], 0);
              avg = sum / res.values.length;
            }
          }

          metrics.metrics.push({
            msg: res.msg,
            type: res.type,
            siteId: res.siteId || siteId || "",
            nodeId: "",
            success: res.success,
            value: formatKPIValue(res.type, avg),
          });
        }
      } catch (error) {
        logger.error(
          `Error processing site metrics for site ${siteId}: ${error}`
        );
        for (const key of metricKeys) {
          metrics.metrics.push({
            msg: "Error",
            type: key,
            siteId: siteId || "",
            nodeId: "",
            success: false,
            value: 0,
          });
        }
      }

      if (Array.isArray(nodeIds) && nodeIds.length > 0) {
        for (const nodeId of nodeIds) {
          try {
            const nodeResults = await Promise.all(
              metricKeys.map(async key => {
                try {
                  return await getNodeMetricRange(baseURL, key, {
                    ...data,
                    siteId,
                    nodeId,
                  });
                } catch (error) {
                  logger.error(
                    `Error processing node metric ${key} for node ${nodeId} in site ${siteId}: ${error}`
                  );
                  return {
                    msg: "Error",
                    type: key,
                    siteId: siteId || "",
                    nodeId: nodeId || "",
                    success: false,
                    values: [],
                  };
                }
              })
            );

            for (const res of nodeResults) {
              let avg = 0;

              if (Array.isArray(res.values)) {
                res.values = res.values.filter(value => value[1] !== 0);
                if (res.values.length > 0) {
                  if (
                    res.values.length === 1 ||
                    res.type === "unit_uptime" ||
                    res.type === "network_uptime"
                  ) {
                    avg = res.values[res.values.length - 1][1];
                  } else {
                    const sum = res.values.reduce(
                      (acc, val) => acc + val[1],
                      0
                    );
                    avg = sum / res.values.length;
                  }
                }
              }

              metrics.metrics.push({
                msg: res.msg,
                type: res.type,
                siteId: siteId || "",
                nodeId: nodeId || "",
                success: res.success,
                value: formatKPIValue(res.type, avg),
              });
            }
          } catch (error) {
            logger.error(
              `Error processing node metrics for node ${nodeId} in site ${siteId}: ${error}`
            );

            for (const key of metricKeys) {
              metrics.metrics.push({
                msg: "Error",
                type: key,
                siteId: siteId || "",
                nodeId: nodeId || "",
                success: false,
                value: 0,
              });
            }
          }
        }
      }
    }

    if (withSubscription && metrics.metrics.length > 0) {
      const baseTopic = `stat-${data.orgName}-${userId}-${type}-${from}`;

      for (const siteId of siteIds) {
        for (const metricKey of metricKeys) {
          const siteUrlParams = new URLSearchParams();
          siteUrlParams.append("interval", data.step.toString());
          siteUrlParams.append("metric", metricKey);
          siteUrlParams.append("site", siteId);

          const siteMetricUrl = `${wsUrl}/v1/live/metrics?${siteUrlParams.toString()}`;

          const siteWorker = new Worker(WS_THREAD, {
            workerData: { topic: baseTopic, url: siteMetricUrl },
          });

          siteWorker.on("message", (_data: any) => {
            if (!_data.isError) {
              try {
                const res = JSON.parse(_data.data);

                if (res?.data?.result && res.data.result.length > 0) {
                  res.data.result.forEach((result: any) => {
                    if (
                      result &&
                      result.metric &&
                      result.value &&
                      result.value.length > 0
                    ) {
                      const resultSiteId =
                        result.metric.site || result.metric.instance || siteId;

                      pubSub.publish(baseTopic, {
                        success: true,
                        msg: "success",
                        type: res.Name,
                        siteId: resultSiteId,
                        nodeId: "",
                        value: [
                          Math.floor(result.value[0]) * 1000,
                          formatKPIValue(res.Name, result.value[1]),
                        ],
                      });
                    }
                  });
                }
              } catch (error) {
                logger.error(
                  `Failed to parse WebSocket message for ${siteId}/${metricKey}: ${error}`
                );
              }
            }
          });

          siteWorker.on("exit", async (code: any) => {
            logger.info(
              `WS_THREAD for site ${siteId}, metric ${metricKey} exited with code [${code}] for ${baseTopic}`
            );
          });

          if (Array.isArray(nodeIds) && nodeIds.length > 0) {
            for (const nodeId of nodeIds) {
              const nodeUrlParams = new URLSearchParams();
              nodeUrlParams.append("interval", data.step.toString());
              nodeUrlParams.append("metric", metricKey);
              nodeUrlParams.append("node", nodeId);

              const nodeMetricUrl = `${wsUrl}/v1/live/metrics?${nodeUrlParams.toString()}`;

              const nodeWorker = new Worker(WS_THREAD, {
                workerData: { topic: baseTopic, url: nodeMetricUrl },
              });

              nodeWorker.on("message", (_data: any) => {
                if (!_data.isError) {
                  try {
                    const res = JSON.parse(_data.data);

                    if (res?.data?.result && res.data.result.length > 0) {
                      res.data.result.forEach((result: any) => {
                        if (
                          result &&
                          result.metric &&
                          result.value &&
                          result.value.length > 0
                        ) {
                          const resultNodeId =
                            result.metric.node ||
                            result.metric.nodeid ||
                            nodeId;

                          pubSub.publish(baseTopic, {
                            success: true,
                            msg: "success",
                            type: res.Name,
                            siteId: siteId,
                            nodeId: resultNodeId,
                            value: [
                              Math.floor(result.value[0]) * 1000,
                              formatKPIValue(res.Name, result.value[1]),
                            ],
                          });
                        }
                      });
                    }
                  } catch (error) {
                    logger.error(
                      `Failed to parse WebSocket message for node ${nodeId}/${metricKey}: ${error}`
                    );
                  }
                }
              });

              nodeWorker.on("exit", async (code: any) => {
                logger.info(
                  `WS_THREAD for node ${nodeId}, metric ${metricKey} exited with code [${code}] for ${baseTopic}`
                );
              });
            }
          }
        }
      }
    }

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

    const wsUrl = wsUrlResolver(baseURL);

    const { type, from, userId, withSubscription, nodeId } = data;
    if (from === 0) throw new Error("Argument 'from' can't be zero.");

    const metricsKey = getGraphsKeyByType(type);
    const metrics: MetricsRes = { metrics: [] };

    if (metricsKey.length > 0) {
      const metricPromises = metricsKey.map(
        async key => await getNodeMetricRange(baseURL, key, { ...data })
      );

      metrics.metrics = await Promise.all(metricPromises);
    }

    if (withSubscription && metrics.metrics.length > 0) {
      const workerData = {
        topic: `${userId}/${type}/${from}`,
        url: `${wsUrl}/v1/live/metrics?interval=${
          data.step
        }&metric=${metricsKey.join(",")}&node=${nodeId}`,
      };

      const worker = new Worker(WS_THREAD, {
        workerData,
      });

      worker.on("message", (_data: any) => {
        if (!_data.isError) {
          try {
            const res = JSON.parse(_data.data);
            const result = res.data.result[0];
            if (result && result.metric && result.value.length > 0) {
              pubSub.publish(workerData.topic, {
                type: res.Name,
                success: true,
                msg: "success",
                nodeId: nodeId,
                value: [
                  Math.floor(result.value[0]) * 1000,
                  parseFloat(Number(result.value[1] || 0).toFixed(2)),
                ],
              });
            }
          } catch (error) {
            logger.error(`Failed to parse WebSocket message: ${error}`);
          }
        }
      });
      worker.on("exit", async (code: any) => {
        await store.close();
        logger.info(
          `WS_THREAD exited with code [${code}] for ${userId}/${type}/${from}`
        );
      });
    }

    return metrics;
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
      logger.error(`Error getting base URL for site metrics: ${baseURL}`);
      return { metrics: [] };
    }
    logger.info(`Using metrics base URL for site: ${baseURL}`);

    const wsUrl = wsUrlResolver(baseURL);

    const { type, from, userId, withSubscription, siteId } = data;
    if (from === 0) throw new Error("Argument 'from' can't be zero.");

    const metricsKey = getGraphsKeyByType(type);
    const metrics: MetricsRes = { metrics: [] };

    if (metricsKey.length > 0) {
      const metricPromises = metricsKey.map(async key => {
        const result = await getSiteMetricRange(baseURL, key, { ...data });
        return result;
      });
      metrics.metrics = await Promise.all(metricPromises);
    }

    if (withSubscription && metrics.metrics.length > 0) {
      const workerData = {
        topic: `${userId}/${type}/${from}`,
        url: `${wsUrl}/v1/live/metrics?interval=${
          data.step
        }&metric=${metricsKey.join(",")}&site=${siteId}`,
      };

      const worker = new Worker(WS_THREAD, {
        workerData,
      });

      worker.on("message", (_data: any) => {
        if (!_data.isError) {
          try {
            const res = JSON.parse(_data.data);
            const result = res.data.result[0];
            if (result && result.metric && result.value.length > 0) {
              pubSub.publish(workerData.topic, {
                type: res.Name,
                success: true,
                msg: "success",
                siteId: siteId,
                value: [
                  Math.floor(result.value[0]) * 1000,
                  parseFloat(Number(result.value[1] || 0).toFixed(2)),
                ],
              });
            }
          } catch (error) {
            logger.error(
              `Failed to parse WebSocket message for site: ${error}`
            );
          }
        }
      });

      worker.on("exit", async (code: any) => {
        await store.close();
        logger.info(
          `WS_THREAD exited with code [${code}] for ${userId}/${type}/${from}`
        );
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
