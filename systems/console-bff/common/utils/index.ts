/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { exec } from "child_process";
import { readFile } from "fs";
import { RootDatabase } from "lmdb";

import InitAPI from "../../init/datasource/init_api";
import { MetricRes, MetricsRes } from "../../subscriptions/resolvers/types";
import {
  GRAPHS_TYPE,
  NOTIFICATION_SCOPE,
  ROLE_TYPE,
  STATS_TYPE,
} from "../enums";
import { HTTP401Error, Messages } from "../errors";
import { logger } from "../logger";
import { Meta, ResponseObj, THeaders } from "../types";
import { RoleToNotificationScopes } from "../utils/roleToNotificationScope";

const getTimestampCount = (count: string) =>
  parseInt((Date.now() / 1000).toString()) + "-" + count;

const parseHeaders = (reqHeader: any): THeaders => {
  const headers: THeaders = {
    auth: {
      Authorization: "",
      Cookie: "",
    },
    token: "",
    orgId: "",
    userId: "",
    orgName: "",
  };
  if (reqHeader.get("introspection") === "true") return headers;
  if (reqHeader.get("x-session-token") ?? reqHeader.get("cookie")) {
    if (reqHeader.get("x-session-token")) {
      headers.auth.Authorization = reqHeader["x-session-token"] as string;
    } else {
      const cookie: string = reqHeader.get("cookie");
      const cookies = cookie.split(";");
      const session: string =
        cookies.find(item => (item.includes("ukama_session") ? item : "")) ??
        "";
      headers.auth.Cookie = session;
      const t =
        cookies.find(item =>
          !item.includes("csrf_token") && item.includes("token") ? item : ""
        ) ?? "";

      if (t !== "") {
        headers.token = t.replace("token=", "");
      } else {
        throw new HTTP401Error(Messages.TOKEN_HEADER_NOT_FOUND);
      }
    }
  } else {
    throw new HTTP401Error(Messages.HEADER_ERR_AUTH);
  }
  return headers;
};

const parseToken = (token: string, get: "orgId" | "orgName" | "userId") => {
  const headers: THeaders = {
    auth: {
      Authorization: "",
      Cookie: "",
    },
    token: "",
    orgId: "",
    userId: "",
    orgName: "",
  };

  if (token) {
    const decoded = Buffer.from(token, "base64").toString("utf-8");
    const headersStr = decoded.split(";");
    if (headersStr.length < 3) throw new HTTP401Error(Messages.HEADER_ERR_AUTH);
    headers.orgId = headersStr[0];
    headers.orgName = headersStr[1];
    headers.userId = headersStr[2];
    headers.token = token;
    return headers[get];
  }
};

const parseGatewayHeaders = (reqHeader: any): THeaders => {
  return {
    auth: {
      Authorization: reqHeader["x-session-token"] ?? "",
      Cookie: reqHeader["cookie"] ?? "",
    },
    token: reqHeader["token"] ?? "",
    orgId: parseToken(reqHeader["token"], "orgId") ?? "",
    userId: parseToken(reqHeader["token"], "userId") ?? "",
    orgName: parseToken(reqHeader["token"], "orgName") ?? "",
  };
};

const getStripeIdByUserId = (uid: string): string => {
  return uid === "d0a36c51-6a66-4187-b786-72a9e09bf7a4"
    ? "cus_MFTZKUVOGtI2fU"
    : "";
};

const getPaginatedOutput = (
  page: number,
  pageSize: number,
  count: number
): Meta => {
  return {
    count,
    page: page ? page : 1,
    size: pageSize ? pageSize : count,
    pages: pageSize ? Math.ceil(count / pageSize) : 1,
  };
};

const getGraphsKeyByType = (type: string): string[] => {
  // TODO: NEED TO UPDATE KPI KEYS
  switch (type) {
    case GRAPHS_TYPE.HOME:
    case STATS_TYPE.HOME:
      return [
        "package_sales",
        "data_usage",
        "node_active_subscribers",
        "network_uptime",
      ];
    case GRAPHS_TYPE.NODE_HEALTH:
      return ["unit_uptime", "unit_health", "node_load"];
    case GRAPHS_TYPE.SUBSCRIBERS:
      return ["subscribers_active"];
    case STATS_TYPE.OVERVIEW:
      return ["unit_uptime", "unit_health", "node_load", "subscribers_active"];
    case GRAPHS_TYPE.NETWORK_CELLULAR:
      return ["cellular_uplink", "cellular_downlink"];
    case GRAPHS_TYPE.NETWORK_BACKHAUL:
      return ["backhaul_uplink", "backhaul_downlink", "backhaul_latency"];
    case STATS_TYPE.NETWORK:
      return [
        "cellular_uplink",
        "cellular_downlink",
        "backhaul_uplink",
        "backhaul_downlink",
        "backhaul_latency",
      ];
    case STATS_TYPE.RESOURCES:
    case GRAPHS_TYPE.RESOURCES:
      return ["hwd_load", "memory_usage", "cpu_usage", "disk_usage"];
    case STATS_TYPE.RADIO:
    case GRAPHS_TYPE.RADIO:
      return ["txpower"];
    case STATS_TYPE.ALL_NODE:
      return [
        "unit_uptime",
        "unit_health",
        "node_load",
        "subscribers_active",
        "cellular_uplink",
        "cellular_downlink",
        "backhaul_uplink",
        "backhaul_downlink",
        "backhaul_latency",
        "hwd_load",
        "memory_usage",
        "cpu_usage",
        "disk_usage",
        "txpower",
      ];
    case GRAPHS_TYPE.BATTERY:
      return ["battery_charge_percentage"];
    case GRAPHS_TYPE.SOLAR:
      return [
        "solar_panel_voltage",
        "solar_panel_current",
        "solar_panel_power",
      ];
    case GRAPHS_TYPE.CONTROLLER:
      return [
        "solar_panel_voltage",
        "solar_panel_current",
        "solar_panel_power",
        "battery_charge_percentage",
      ];
    case GRAPHS_TYPE.MAIN_BACKHAUL:
      return ["main_backhaul_latency", "backhaul_speed"];
    case GRAPHS_TYPE.DATA_USAGE:
      return ["data_usage"];
    case GRAPHS_TYPE.SWITCH:
      return [
        "backhaul_switch_port_status",
        "backhaul_switch_port_speed",
        "backhaul_switch_port_power",
        "solar_switch_port_status",
        "solar_switch_port_speed",
        "solar_switch_port_power",
        "node_switch_port_status",
        "node_switch_port_speed",
        "node_switch_port_power",
      ];
    case GRAPHS_TYPE.SITE:
      return [
        "site_uptime_seconds",
        "unit_uptime",
        "solar_panel_voltage",
        "solar_panel_current",
        "site_uptime_percentage",
        "solar_panel_power",
        "battery_charge_percentage",
        "main_backhaul_latency",
        "backhaul_speed",
        "backhaul_switch_port_status",
        "backhaul_switch_port_speed",
        "backhaul_switch_port_power",
        "solar_switch_port_status",
        "solar_switch_port_speed",
        "solar_switch_port_power",
        "node_switch_port_status",
        "node_switch_port_speed",
        "node_switch_port_power",
        "node_active_subscribers",
      ];
    default:
      return [];
  }
};

const findProcessNKill = (port: string): Promise<boolean> => {
  return new Promise((resolve, reject) => {
    const command = `lsof -i tcp:${port} | awk 'NR>1 {print $2}'`;

    exec(command, (err, stdout) => {
      if (err) {
        reject(new Error(`Failed to execute command: ${err.message}`));
        return;
      }
      if (stdout) {
        const pid = stdout.replace(/\n/g, "");

        if (!pid) {
          reject(new Error("PID not found."));
          return;
        }
        killProcess(pid)
          .then(() => resolve(true))
          .catch(error => reject(error));
      } else {
        resolve(true);
      }
    });
  });
};

const killProcess = (pid: string): Promise<void> => {
  return new Promise((resolve, reject) => {
    const command = `kill -9 ${pid}`;

    exec(command, err => {
      if (err) {
        reject(new Error(`Error killing process ${pid}: ${err.message}`));
      } else {
        logger.info(`Process ${pid} killed.`);
        resolve();
      }
    });
  });
};

const getSystemNameByService = (service: string): string => {
  switch (service) {
    case "org":
    case "user":
      return "nucleus";
    case "network":
    case "member":
    case "site":
    case "invitation":
    case "node":
      return "registry";
    case "package":
    case "rate":
      return "dataplan";
    case "sim":
    case "subscriber":
      return "subscriber";
    case "notification":
      return "notification";
    case "init":
      return "init";
    case "billing":
      return "billing";
    case "report":
      return "report";
    case "payments":
      return "payments";
    case "metrics":
      return "metrics";
    case "planning-tool":
      return "planning";
    case "nodeState":
    case "controller":
      return "node";
    default:
      return "";
  }
};

const getBaseURL = async (
  serviceName: string,
  orgName: string,
  store?: RootDatabase
): Promise<ResponseObj> => {
  const sysName = getSystemNameByService(serviceName);
  logger.info(`${store?.get("org")}`);

  const initAPI = new InitAPI();
  if (orgName && sysName) {
    try {
      const intRes = await initAPI.getSystem(orgName, sysName);
      const url = intRes.url
        ? intRes.url
        : `http://${intRes.ip}:${intRes.port}`;
      return {
        status: 200,
        message: url,
      };
    } catch (e) {
      logger.error(`Error getting base URL for ${orgName}-${sysName}: ${e}`);
    }
  }
  return {
    status: 500,
    message: "Unable to reach system",
  };
};

const csvToBase64 = (filePath: string) => {
  readFile(filePath, (err, data) => {
    if (err) {
      logger.error("Error reading file: ", err);
      return;
    }
    return data.toString("base64");
  });
};

const getRoleType = (userRole: string): ROLE_TYPE => {
  return Object.values(ROLE_TYPE).includes(userRole as ROLE_TYPE)
    ? (userRole as ROLE_TYPE)
    : ROLE_TYPE.ROLE_INVALID;
};

const getScopesByRole = (userRole: string): Array<NOTIFICATION_SCOPE> => {
  const roleType = getRoleType(userRole);
  return RoleToNotificationScopes[roleType] ?? [];
};

export const wsUrlResolver = (url: string): string => {
  if (url?.startsWith("wss://") || url?.startsWith("ws://")) {
    return url;
  } else if (url?.includes("https://")) {
    return url.replace("https://", "wss://");
  } else if (url?.startsWith("http://")) {
    return url.replace("http://", "ws://");
  }
  return url;
};

export const formatKPIValue = (type: string, value: any) => {
  if (value === "NaN") return 0;
  switch (type) {
    case "backhaul_latency":
      return Math.floor(Number(value || 0));
    case "subscribers_active":
    case "node_active_subscribers":
      return Math.floor(parseFloat(value || "0"));
    default:
      return parseFloat(Number(value || 0).toFixed(2));
  }
};

const epochToISOString = (epoch: number): string => {
  const date = new Date(epoch * 1000);
  return date.toISOString().replace(/\.\d{3}Z$/, "Z");
};

export const transformMetricsArray = (
  metricsArray: MetricsRes[]
): MetricsRes => {
  const allMetrics: MetricRes[] = metricsArray.flatMap(item => item.metrics);

  return {
    metrics: allMetrics,
  };
};
const handleMetricWSMessage = (
  data: any,
  topic: string,
  pubSub: any,
  defaultSiteId: string = "",
  defaultNodeId: string = ""
) => {
  if (!data.isError) {
    try {
      const res = JSON.parse(data.data);

      if (res?.data?.result && res.data.result.length > 0) {
        res.data.result.forEach((result: any) => {
          if (
            result &&
            result.metric &&
            result.value &&
            result.value.length > 0
          ) {
            const resultSiteId =
              result.metric.site || result.metric.instance || defaultSiteId;
            const resultNodeId =
              result.metric.node || result.metric.nodeid || defaultNodeId;

            pubSub.publish(topic, {
              success: true,
              msg: "success",
              type: res.Name,
              siteId: resultSiteId,
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
      logger.error(`Failed to parse WebSocket message: ${error}`);
    }
  }
};

const isMetricValidNetworkCheck = (
  arg: string,
  res: string,
  op: string
): boolean => {
  if (op !== "sum" && arg && arg !== res) return true;
  return false;
};

export {
  csvToBase64,
  epochToISOString,
  findProcessNKill,
  getBaseURL,
  getGraphsKeyByType,
  getPaginatedOutput,
  getScopesByRole,
  getStripeIdByUserId,
  getSystemNameByService,
  getTimestampCount,
  handleMetricWSMessage,
  isMetricValidNetworkCheck,
  killProcess,
  parseGatewayHeaders,
  parseHeaders,
  parseToken,
};
