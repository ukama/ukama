/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import "dotenv/config";

export const VERSION = process.env.VERSION ?? "v1";

// API GWs
export const PLANNING_API_URL = process.env.PLANNING_API_URL;
export const METRIC_API_GW = process.env.METRIC_API_GW ?? "";
export const NOTIFICATION_API_GW = process.env.NOTIFICATION_API_GW ?? "";
export const NOTIFICATION_API_GW_WS = process.env.NOTIFICATION_API_GW_WS ?? "";
export const METRIC_API_GW_SOCKET = process.env.METRIC_API_GW_SOCKET ?? "";
export const REGISTRY_API_GW = process.env.REGISTRY_API_GW ?? "";
export const SUBSCRIBER_API_GW = process.env.SUBSCRIBER_API_GW ?? "";
export const NUCLEUS_API_GW = process.env.NUCLEUS_API_GW ?? "";
export const DATA_API_GW = process.env.DATA_API_GW ?? "";
export const INIT_API_GW = process.env.INIT_API_GW ?? "";
export const BILLING_API_GW = process.env.BILLING_API_GW ?? "";
export const INVENTORY_API_GW = process.env.INVENTORY_API_GW ?? "";

// FRONTEND URLS
export const AUTH_APP_URL = process.env.AUTH_APP_URL ?? "";
export const PLAYGROUND_URL = process.env.PLAYGROUND_URL ?? "";
export const CONSOLE_APP_URL = process.env.CONSOLE_APP_URL ?? "";

// UTILS
export const BASE_DOMAIN = process.env.BASE_DOMAIN ?? "ukama.com";
export const COMMUNITY_ORG_NAME = process.env.COMMUNITY_ORG_NAME ?? "ukama";
export const ENCRYPTION_KEY = process.env.ENCRYPTION_KEY ?? "";
export const PLANNING_TOOL_DB = process.env.PLANNING_TOOL_DB ?? "";
export const AUTH_URL = process.env.AUTH_URL ?? "";
export const STORAGE_KEY = process.env.STORAGE_KEY ?? "UKAMA_STORAGE_KEY";
export const PLANNING_BUCKET = process.env.BUCKET_NAME;
export const STRIP_SK = process.env.STRIP_SK ?? "";
export const METRIC_PROMETHEUS = process.env.METRIC_PROMETHEUS ?? "";
export const BFF_REDIS = process.env.BFF_REDIS ?? "redis://localhost:6379";

// PORTS
export const GATEWAY_PORT = parseInt(process.env.GATEWAY_PORT ?? "8080");
export const SUBSCRIPTIONS_PORT = parseInt(
  process.env.SUBSCRIPTIONS_PORT ?? "8081"
);
export const PLANNING_SERVICE_PORT = parseInt(
  process.env.PLANNING_SERVICE_PORT ?? "5042"
);
const NODE_PORT = parseInt(process.env.NODE_PORT ?? "5043");
const USER_PORT = parseInt(process.env.USER_PORT ?? "5044");
const PACKAGE_PORT = parseInt(process.env.PACKAGE_PORT ?? "5045");
const RATE_PORT = parseInt(process.env.RATE_PORT ?? "5046");
const ORG_PORT = parseInt(process.env.ORG_PORT ?? "5047");
const NETWORK_PORT = parseInt(process.env.NETWORK_PORT ?? "5048");
export const BILLING_PORT = parseInt(process.env.BILLING_PORT ?? "5051");
const SIM_PORT = parseInt(process.env.SIM_PORT ?? "5052");
const INVITATION_PORT = parseInt(process.env.INVITATION_PORT ?? "5053");
const MEMBER_PORT = parseInt(process.env.MEMBER_PORT ?? "5054");
const INIT_PORT = parseInt(process.env.INIT_PORT ?? "5055");
const SUBSCRIBER_PORT = parseInt(process.env.SUBSCRIBER_PORT ?? "5056");
const NOTIFICATION_PORT = parseInt(process.env.NOTIFICATION_PORT ?? "5057");
export const SITE_PORT = parseInt(process.env.SITE_PORT ?? "5058");
export const COMPONENT_INVENTORY_PORT = parseInt(
  process.env.COMPONENT_INVENTORY_PORT ?? "5059"
);
export const SUB_GRAPHS = {
  org: {
    name: "org",
    port: ORG_PORT,
    url: `http://localhost:${ORG_PORT}`,
    isPingedSuccess: false,
  },
  node: {
    name: "node",
    port: NODE_PORT,
    url: `http://localhost:${NODE_PORT}`,
    isPingedSuccess: false,
  },
  user: {
    name: "user",
    port: USER_PORT,
    url: `http://localhost:${USER_PORT}`,
    isPingedSuccess: false,
  },
  network: {
    name: "network",
    port: NETWORK_PORT,
    url: `http://localhost:${NETWORK_PORT}`,
    isPingedSuccess: false,
  },
  component: {
    name: "component",
    port: COMPONENT_INVENTORY_PORT,
    url: `http://localhost:${COMPONENT_INVENTORY_PORT}`,
    isPingedSuccess: false,
  },
  subscriber: {
    name: "subscriber",
    port: SUBSCRIBER_PORT,
    url: `http://localhost:${SUBSCRIBER_PORT}`,
    isPingedSuccess: false,
  },
  sim: {
    name: "sim",
    port: SIM_PORT,
    url: `http://localhost:${SIM_PORT}`,
    isPingedSuccess: false,
  },
  site: {
    name: "site",
    port: SITE_PORT,
    url: `http://localhost:${SITE_PORT}`,
    isPingedSuccess: false,
  },
  package: {
    name: "package",
    port: PACKAGE_PORT,
    url: `http://localhost:${PACKAGE_PORT}`,
    isPingedSuccess: false,
  },
  rate: {
    name: "rate",
    port: RATE_PORT,
    url: `http://localhost:${RATE_PORT}`,
    isPingedSuccess: false,
  },
  invitation: {
    name: "invitation",
    port: INVITATION_PORT,
    url: `http://localhost:${INVITATION_PORT}`,
    isPingedSuccess: false,
  },
  member: {
    name: "member",
    port: MEMBER_PORT,
    url: `http://localhost:${MEMBER_PORT}`,
    isPingedSuccess: false,
  },
  // {
  //   name: "planning",
  //   url: `http://localhost:${PLANNING_SERVICE_PORT}`,
  //   isPingedSuccess: false,
  // },
  init: {
    name: "init",
    port: INIT_PORT,
    url: `http://localhost:${INIT_PORT}`,
    isPingedSuccess: false,
  },
  notification: {
    name: "notification",
    port: NOTIFICATION_PORT,
    url: `http://localhost:${NOTIFICATION_PORT}`,
    isPingedSuccess: false,
  },
};

export const SUB_GRAPH_LIST = Object.entries(SUB_GRAPHS).map(
  ([key, value]) => ({
    type: key,
    ...value,
  })
);
