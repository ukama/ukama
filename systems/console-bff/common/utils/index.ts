/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { GRAPHS_TYPE, NODE_TYPE } from "../enums";
import { HTTP401Error, Messages } from "../errors";
import { Meta, THeaders } from "../types";

const getTimestampCount = (count: string) =>
  parseInt((Date.now() / 1000).toString()) + "-" + count;

const parseHeaders = (reqHeader: any): THeaders => {
  const headers: THeaders = {
    auth: {
      Authorization: "",
      Cookie: "",
    },
    orgId: "",
    userId: "",
    orgName: "",
  };
  if (reqHeader.get("introspection") === "true") return headers;
  if (reqHeader.get("org-id")) {
    headers.orgId = reqHeader.get("org-id") as string;
  } else {
    throw new HTTP401Error(Messages.HEADER_ERR_ORG);
  }
  if (reqHeader.get("user-id")) {
    headers.userId = reqHeader.get("user-id") as string;
  } else {
    throw new HTTP401Error(Messages.HEADER_ERR_USER);
  }
  if (reqHeader.get("org-name")) {
    headers.orgName = reqHeader.get("org-name") as string;
  } else {
    throw new HTTP401Error(Messages.HEADER_ERR_ORG_NAME);
  }

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
    }
  } else {
    throw new HTTP401Error(Messages.HEADER_ERR_AUTH);
  }
  return headers;
};

const parseGatewayHeaders = (reqHeader: any): THeaders => {
  return {
    auth: {
      Authorization: reqHeader["x-session-token"] ?? "",
      Cookie: reqHeader["cookie"] ?? "",
    },
    orgId: reqHeader["orgid"] ?? "",
    userId: reqHeader["userid"] ?? "",
    orgName: reqHeader["orgname"] ?? "",
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

export {
  getGraphsKeyByType,
  getPaginatedOutput,
  getStripeIdByUserId,
  getTimestampCount,
  parseGatewayHeaders,
  parseHeaders,
};
