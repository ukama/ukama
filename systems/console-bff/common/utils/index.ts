/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { exec } from "child_process";
import { readFile } from "fs";

import InitAPI from "../../init/datasource/init_api";
import { GRAPHS_TYPE, NODE_TYPE } from "../enums";
import { HTTP401Error, Messages } from "../errors";
import { logger } from "../logger";
import { Meta, ResponseObj, THeaders } from "../types";

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
        throw new HTTP401Error(Messages.HEADER_ERR_AUTH);
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

const getGraphsKeyByType = (type: string, nodeId: string): string[] => {
  switch (type) {
    case GRAPHS_TYPE.NODE_HEALTH:
      if (nodeId.includes(NODE_TYPE.hnode))
        return ["uptime_trx", "temperature_trx", "temperature_rfe"];
      else if (nodeId.includes(NODE_TYPE.anode))
        return ["temperature_ctl", "temperature_rfe"];
      else return ["temperature_trx", "temperature_com"];
    case GRAPHS_TYPE.NETWORK:
      if (!nodeId.includes(NODE_TYPE.anode))
        return ["rrc", "rlc", "erab", "throughputuplink", "throughputdownlink"];
      else return [];
    case GRAPHS_TYPE.RESOURCES:
      if (nodeId.includes(NODE_TYPE.hnode))
        return ["cpu_trx_usage", "memory_trx_used", "disk_trx_used"];
      else if (nodeId.includes(NODE_TYPE.anode))
        return ["cpu_ctl_used", "disk_ctl_used", "memory_ctl_used"];
      else
        return [
          "power_level",
          "cpu_trx_usage",
          "cpu_com_usage",
          "disk_trx_used",
          "disk_com_used",
          "memory_trx_used",
          "memory_com_used",
        ];
    case GRAPHS_TYPE.RADIO:
      if (nodeId.includes(NODE_TYPE.hnode))
        return ["tx_power", "rx_power", "pa_power"];
      else return [];
    case GRAPHS_TYPE.SUBSCRIBERS:
      if (nodeId.includes(NODE_TYPE.hnode))
        return ["subscribers_active", "subscribers_attached"];
      else if (nodeId.includes(NODE_TYPE.anode))
        return ["temperature_ctl", "temperature_rfe"];
      else return ["subscribers_active", "subscribers_attached"];
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
    case "planning-tool":
      return "planning";
    default:
      return "";
  }
};

const getBaseURL = async (
  serviceName: string,
  orgName: string,
  redisClient: any
): Promise<ResponseObj> => {
  const sysName = getSystemNameByService(serviceName);
  if (redisClient) {
    const redisBaseURL = await redisClient.get(`${sysName}-${orgName}`);
    if (redisBaseURL)
      return {
        status: 200,
        message: redisBaseURL,
      };
  }

  const initAPI = new InitAPI();
  if (orgName && sysName) {
    const intRes = await initAPI.getSystem(orgName, sysName);
    if (redisClient) await redisClient.set(`${sysName}-${orgName}`, intRes.url);
    return {
      status: 200,
      message: intRes.url ? intRes.url : `http://${intRes.ip}:${intRes.port}`,
    };
  } else {
    return {
      status: 500,
      message: "Unable to reach system",
    };
  }
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

export {
  csvToBase64,
  findProcessNKill,
  getBaseURL,
  getGraphsKeyByType,
  getPaginatedOutput,
  getStripeIdByUserId,
  getSystemNameByService,
  getTimestampCount,
  killProcess,
  parseGatewayHeaders,
  parseHeaders,
  parseToken,
};
