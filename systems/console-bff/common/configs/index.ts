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
export const NUCLEUS_API_GW = process.env.NUCLEUS_API_GW ?? "";
export const INIT_API_GW = process.env.INIT_API_GW ?? "";
export const INVENTORY_API_GW = process.env.INVENTORY_API_GW ?? "";

// FRONTEND URLS
export const AUTH_APP_URL = process.env.AUTH_APP_URL ?? "http://localhost:4455";
export const PLAYGROUND_URL =
  process.env.PLAYGROUND_URL ?? "http://localhost:8080";
export const CONSOLE_APP_URL =
  process.env.CONSOLE_APP_URL ?? "http://localhost:3000";

// UTILS
export const BASE_DOMAIN = process.env.BASE_DOMAIN ?? "ukama.com";
export const COMMUNITY_ORG_NAME = process.env.COMMUNITY_ORG_NAME ?? "ukama";
export const ENCRYPTION_KEY = process.env.ENCRYPTION_KEY ?? "";
export const PLANNING_TOOL_DB = process.env.PLANNING_TOOL_DB ?? "";
export const AUTH_URL = process.env.AUTH_URL ?? "";
export const STORAGE_KEY = process.env.STORAGE_KEY ?? "UKAMA_STORAGE_KEY";
export const STRIP_SK = process.env.STRIP_SK ?? "";

// PORTS
export const GATEWAY_PORT = parseInt(process.env.GATEWAY_PORT ?? "8080");
export const SUBSCRIPTIONS_PORT = parseInt(
  process.env.SUBSCRIPTIONS_PORT ?? "8081"
);
export const PLANNING_SERVICE_PORT = parseInt(
  process.env.PLANNING_SERVICE_PORT ?? "5041"
);
const ORG_PORT = parseInt(process.env.ORG_PORT ?? "5042");
const USER_PORT = parseInt(process.env.USER_PORT ?? "5043");
const INIT_PORT = parseInt(process.env.INIT_PORT ?? "5044");
const PACKAGE_PORT = parseInt(process.env.PACKAGE_PORT ?? "5045");
const RATE_PORT = parseInt(process.env.RATE_PORT ?? "5046");
const NETWORK_PORT = parseInt(process.env.NETWORK_PORT ?? "5047");
const SITE_PORT = parseInt(process.env.SITE_PORT ?? "5048");
const INVITATION_PORT = parseInt(process.env.INVITATION_PORT ?? "5049");
const MEMBER_PORT = parseInt(process.env.MEMBER_PORT ?? "5050");
const NODE_PORT = parseInt(process.env.NODE_PORT ?? "5051");
const SUBSCRIBER_PORT = parseInt(process.env.SUBSCRIBER_PORT ?? "5052");
const SIM_PORT = parseInt(process.env.SIM_PORT ?? "5053");
const NOTIFICATION_PORT = parseInt(process.env.NOTIFICATION_PORT ?? "5054");
const CONTROLLER_PORT = parseInt(process.env.CONTROLLER_PORT ?? "5058");

export const BILLING_PORT = parseInt(process.env.BILLING_PORT ?? "5055");
export const COMPONENT_INVENTORY_PORT = parseInt(
  process.env.COMPONENT_INVENTORY_PORT ?? "5056"
);
export const METRIC_PORT = parseInt(process.env.METRIC_PORT ?? "5057");
export const SUB_GRAPHS = {
  metric: {
    name: "metric",
    port: METRIC_PORT,
    url: `http://localhost:${METRIC_PORT}`,
    isPingedSuccess: false,
  },
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
  controller: {
    name: "controller",
    port: CONTROLLER_PORT,
    url: `http://localhost:${CONTROLLER_PORT}`,
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
