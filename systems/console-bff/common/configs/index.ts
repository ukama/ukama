/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import "dotenv/config";

export const VERSION = process.env.VERSION || "v1";

// API GWs
export const PLANNING_API_URL = process.env.PLANNING_API_URL;
export const METRIC_API_GW = process.env.METRIC_API_GW || "";
export const METRIC_API_GW_SOCKET = process.env.METRIC_API_GW_SOCKET || "";
export const REGISTRY_API_GW = process.env.REGISTRY_API_GW || "";
export const SUBSCRIBER_API_GW = process.env.SUBSCRIBER_API_GW || "";
export const NUCLEUS_API_GW = process.env.NUCLEUS_API_GW || "";
export const DATA_API_GW = process.env.DATA_API_GW || "";
export const BILLING_API_GW = process.env.BILLING_API_GW || "";

// FRONTEND URLS
export const AUTH_APP_URL = process.env.AUTH_APP_URL || "";
export const PLAYGROUND_URL = process.env.PLAYGROUND_URL || "";
export const CONSOLE_APP_URL = process.env.CONSOLE_APP_URL || "";

// UTILS
export const AUTH_URL = process.env.AUTH_URL || "";
export const STORAGE_KEY = process.env.STORAGE_KEY || "";
export const PLANNING_BUCKET = process.env.BUCKET_NAME;
export const STRIP_SK = process.env.STRIP_SK || "";
export const METRIC_PROMETHEUS = process.env.METRIC_PROMETHEUS || "";

// PORTS
export const GATEWAY_PORT = parseInt(process.env.GATEWAY_PORT || "8000");
export const PLANNING_SERVICE_PORT = parseInt(
  process.env.PLANNING_SERVICE_PORT || "5041"
);
export const METRICS_PORT = parseInt(process.env.METRICS_PORT || "5042");
export const NODE_PORT = parseInt(process.env.NODE_PORT || "5043");
export const USER_PORT = parseInt(process.env.USER_PORT || "5044");
export const PACKAGE_PORT = parseInt(process.env.PACKAGE_PORT || "5045");
export const RATE_PORT = parseInt(process.env.RATE_PORT || "5046");
export const ORG_PORT = parseInt(process.env.ORG_PORT || "5047");
export const NETWORK_PORT = parseInt(process.env.NETWORK_PORT || "5048");
export const SUBSCRIBER_PORT = parseInt(process.env.SUBSCRIBER_PORT || "5049");
export const ALERT_PORT = parseInt(process.env.ALERT_PORT || "5050");
export const BILLING_PORT = parseInt(process.env.BILLING_PORT || "5051");
export const SIM_PORT = parseInt(process.env.SIM_PORT || "5052");
export const INVITATION_PORT = parseInt(process.env.INVITATION_PORT || "5053");
export const MEMBER_PORT = parseInt(process.env.MEMBER_PORT || "5054");
