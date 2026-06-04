/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { exec } from "child_process";
import { readFile } from "fs";
import { IncomingHttpHeaders } from "http";
import { RootDatabase } from "lmdb";

import InitAPI from "../../init/datasource/init_api";
import { MetricRes, MetricsRes } from "../../subscriptions/resolvers/types";
import { verifyToken } from "../auth/token";
import { SUB_GRAPHS } from "../configs";
import {
  GRAPHS_TYPE,
  NODE_TYPE,
  NOTIFICATION_SCOPE,
  ROLE_TYPE,
  STATS_TYPE,
} from "../enums";
import { Messages, UnauthenticatedError } from "../errors";
import { logger } from "../logger";
import { Meta, ResponseObj, THeaders } from "../types";
import { RoleToNotificationScopes } from "../utils/roleToNotificationScope";

const getTimestampCount = (count: string) =>
  parseInt((Date.now() / 1000).toString()) + "-" + count;

// Index of the optional `exp` (epoch seconds) claim within a token.
const TOKEN_EXP_INDEX = 10;

/**
 * Verifies a token's HMAC signature and (if present) its expiry, returning
 * the decoded claim segments. Returns null for a missing/forged signature,
 * a malformed payload, or an expired token. Tokens issued before the `exp`
 * claim existed (fewer segments) remain valid for backward compatibility.
 */
const verifyTokenClaims = (token: string): string[] | null => {
  const payload = verifyToken(token);
  if (!payload) return null;

  const claims = Buffer.from(payload, "base64").toString("utf-8").split(";");
  if (claims.length < 3) return null;

  const expRaw = claims[TOKEN_EXP_INDEX];
  if (expRaw) {
    const exp = parseInt(expRaw, 10);
    if (!Number.isNaN(exp) && Math.floor(Date.now() / 1000) >= exp) {
      return null; // expired
    }
  }
  return claims;
};

const parseToken = (
  token: string,
  get: "orgId" | "orgName" | "userId"
): string | undefined => {
  if (!token) return undefined;

  // Reject tokens whose signature is missing/invalid or that have expired,
  // so clients cannot forge or replay org/user claims.
  const claims = verifyTokenClaims(token);
  if (!claims) throw new UnauthenticatedError(Messages.HEADER_ERR_AUTH);

  switch (get) {
    case "orgId":
      return claims[0];
    case "orgName":
      return claims[1];
    case "userId":
      return claims[2];
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

/**
 * Parses auth headers from a plain Express request headers object
 * (used by the gateway's per-request context). Extracts the session
 * cookie and the signed token, throwing 401 when neither a session
 * token header nor a session cookie is present.
 */
const parseExpressHeaders = (reqHeader: IncomingHttpHeaders): THeaders => {
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

  const sessionToken = (reqHeader["x-session-token"] as string) ?? "";
  const cookieHeader = reqHeader["cookie"] ?? "";

  if (!sessionToken && !cookieHeader) {
    throw new UnauthenticatedError(Messages.HEADER_ERR_AUTH);
  }

  if (sessionToken) {
    headers.auth.Authorization = sessionToken;
    return headers;
  }

  const cookies = cookieHeader.split(";").map(c => c.trim());
  headers.auth.Cookie =
    cookies.find(item => item.includes("ukama_session")) ?? "";

  const tokenCookie =
    cookies.find(
      item => !item.includes("csrf_token") && item.startsWith("token=")
    ) ?? "";
  if (!tokenCookie) {
    throw new UnauthenticatedError(Messages.TOKEN_HEADER_NOT_FOUND);
  }

  const token = tokenCookie.replace("token=", "");
  // Verify signature + expiry at the gateway entry point so forged, tampered,
  // or expired tokens are rejected with a clean 401 before any subgraph runs.
  if (!verifyTokenClaims(token)) {
    throw new UnauthenticatedError(Messages.HEADER_ERR_AUTH);
  }
  headers.token = token;
  return headers;
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

const TYPE_KEYS_GROUPS: { types: string[]; keys: string[] }[] = [
  {
    types: [GRAPHS_TYPE.HOME, STATS_TYPE.HOME],
    keys: [
      "package_sales",
      "data_usage",
      "node_active_subscribers",
      "network_uptime",
    ],
  },
  {
    types: [STATS_TYPE.RESOURCES, GRAPHS_TYPE.RESOURCES],
    keys: ["cpu", "memory", "disk"],
  },
  {
    types: [GRAPHS_TYPE.DATA_USAGE],
    keys: ["data_usage"],
  },
];

const TYPE_NODE_KEYS_GROUPS: {
  types: string[];
  nodeKeys: Partial<Record<NODE_TYPE, string[]>>;
}[] = [
  {
    types: [GRAPHS_TYPE.NODE_HEALTH],
    nodeKeys: {
      [NODE_TYPE.tnode]: ["uptime", "cpu_temperature", "memory"],
      [NODE_TYPE.anode]: ["uptime", "fem1_temperature", "fem2_temperature"],
      [NODE_TYPE.cnode]: ["uptime", "memory"],
    },
  },
  {
    types: [GRAPHS_TYPE.SUBSCRIBERS],
    nodeKeys: {
      [NODE_TYPE.tnode]: ["subscribers_active"],
    },
  },
  {
    types: [STATS_TYPE.OVERVIEW],
    nodeKeys: {
      [NODE_TYPE.tnode]: [
        "uptime",
        "cpu_temperature",
        "memory",
        "subscribers_active",
      ],
      [NODE_TYPE.anode]: ["uptime", "fem1_temperature", "fem2_temperature"],
      [NODE_TYPE.cnode]: ["uptime", "memory"],
    },
  },
  {
    types: [GRAPHS_TYPE.NETWORK_CELLULAR],
    nodeKeys: {
      [NODE_TYPE.tnode]: ["cellular_uplink", "cellular_downlink"],
    },
  },
  {
    types: [GRAPHS_TYPE.NETWORK_BACKHAUL],
    nodeKeys: {
      [NODE_TYPE.tnode]: [
        "backhaul_uplink",
        "backhaul_downlink",
        "backhaul_latency",
      ],
    },
  },
  {
    types: [STATS_TYPE.NETWORK],
    nodeKeys: {
      [NODE_TYPE.tnode]: [
        "cellular_uplink",
        "cellular_downlink",
        "backhaul_uplink",
        "backhaul_downlink",
        "backhaul_latency",
      ],
    },
  },
  {
    types: [STATS_TYPE.RADIO, GRAPHS_TYPE.RADIO],
    nodeKeys: {
      [NODE_TYPE.tnode]: ["power"],
      [NODE_TYPE.anode]: ["pa_power", "rx_power", "tx_power"],
    },
  },
  {
    types: [STATS_TYPE.ALL_NODE],
    nodeKeys: {
      [NODE_TYPE.tnode]: [
        "uptime",
        "cpu_temperature",
        "subscribers_active",
        "cellular_uplink",
        "cellular_downlink",
        "backhaul_uplink",
        "backhaul_downlink",
        "backhaul_latency",
        "cpu",
        "memory",
        "disk",
        "power",
      ],
      [NODE_TYPE.anode]: [
        "uptime",
        "fem1_temperature",
        "fem2_temperature",
        "cpu",
        "memory",
        "disk",
        "pa_power",
        "rx_power",
        "tx_power",
      ],
      [NODE_TYPE.cnode]: ["uptime", "cpu", "memory", "disk"],
    },
  },
];

const getInfraGraphKpiKeysByType = (type: string): string[] => {
  switch (type) {
    case GRAPHS_TYPE.BATTERY:
      return ["battery_charge"];
    case GRAPHS_TYPE.SOLAR:
      return [
        "solar_panel_voltage",
        "solar_panel_current",
        "solar_panel_power",
      ];
    case GRAPHS_TYPE.CONTROLLER:
      return ["controller_temperature", "load_current"];
    case GRAPHS_TYPE.MAIN_BACKHAUL:
      return ["backhaul_latency", "backhaul_downlink"];
    case GRAPHS_TYPE.SWITCH:
      return [
        "switch_port_1_speed",
        "switch_port_1_power",
        "switch_port_2_speed",
        "switch_port_2_power",
        "switch_port_3_speed",
        "switch_port_3_power",
        "switch_port_4_speed",
        "switch_port_4_power",
        "switch_port_9_speed",
        "switch_port_9_power",
      ];
    case GRAPHS_TYPE.SITE:
      return [
        "solar_panel_voltage",
        "solar_panel_current",
        "solar_panel_power",
        "controller_temperature",
        "load_current",
        "battery_charge",
        "backhaul_latency",
        "backhaul_downlink",
        "switch_port_1_speed",
        "switch_port_1_power",
        "switch_port_2_speed",
        "switch_port_2_power",
        "switch_port_3_speed",
        "switch_port_3_power",
        "switch_port_4_speed",
        "switch_port_4_power",
        "switch_port_9_speed",
        "switch_port_9_power",
        "node_active_subscribers",
      ];
    default:
      return [];
  }
};

const getGraphsKeyByType = (type: string, nodeType: NODE_TYPE): string[] => {
  const typeGroup = TYPE_KEYS_GROUPS.find(group => group.types.includes(type));
  if (typeGroup) {
    return typeGroup.keys;
  }

  const nodeTypeGroup = TYPE_NODE_KEYS_GROUPS.find(group =>
    group.types.includes(type)
  );
  return nodeTypeGroup?.nodeKeys[nodeType] ?? [];
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
    case "state":
    case "health":
    case "controller":
    case "software":
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
  const isForNodeGw = SUB_GRAPHS[serviceName]?.isForNodeGw ?? false;
  const sysName = getSystemNameByService(serviceName);
  logger.info(`${store?.get("org")}`);

  const initAPI = new InitAPI();
  if (orgName && sysName) {
    try {
      const intRes = await initAPI.getSystem(orgName, sysName);
      if (isForNodeGw) {
        return {
          status: 200,
          message: `http://${intRes.nodeGwIp}:${intRes.nodeGwPort}`,
        };
      } else {
        return {
          status: 200,
          message: intRes.apiGwUrl
            ? intRes.apiGwUrl
            : `http://${intRes.apiGwIp}:${intRes.apiGwPort}`,
        };
      }
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
    case "uptime":
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

const isMetricNetworkCheckFailed = (
  arg: string,
  res: string,
  op: string
): boolean => {
  if (op !== "sum" && arg && arg !== res) return true;
  return false;
};

const getNodeTypeFromId = (id: string): NODE_TYPE => {
  if (id.includes("tnode")) return NODE_TYPE.tnode;
  if (id.includes("anode")) return NODE_TYPE.anode;
  if (id.includes("hnode")) return NODE_TYPE.hnode;
  if (id.includes("cnode")) return NODE_TYPE.cnode;
  return NODE_TYPE.tnode;
};
export {
  csvToBase64,
  epochToISOString,
  findProcessNKill,
  getBaseURL,
  getGraphsKeyByType,
  getInfraGraphKpiKeysByType,
  getNodeTypeFromId,
  getPaginatedOutput,
  getScopesByRole,
  getSystemNameByService,
  getTimestampCount,
  handleMetricWSMessage,
  isMetricNetworkCheckFailed,
  killProcess,
  parseExpressHeaders,
  parseGatewayHeaders,
  parseToken,
};
