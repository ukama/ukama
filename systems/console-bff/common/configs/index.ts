import "dotenv/config";

// API GWs
export const PLANNING_API_URL = process.env.PLANNING_API_URL;
export const METRIC_API_GW = process.env.METRIC_API_GW || "";

// FRONTEND URLS
export const AUTH_APP_URL = process.env.AUTH_APP_URL || "";
export const PLAYGROUND_URL = process.env.PLAYGROUND_URL || "";
export const CONSOLE_APP_URL = process.env.CONSOLE_APP_URL || "";

// UTILS
export const STORAGE_KEY = process.env.STORAGE_KEY || "";
export const PLANNING_BUCKET = process.env.BUCKET_NAME;

// PORTS
export const PLANNING_SERVICE_PORT = parseInt(
  process.env.PLANNING_SERVICE_PORT || "4041"
);
export const GATEWAY_PORT = parseInt(process.env.GATEWAY_PORT || "4000");
export const METRICS_PORT = parseInt(process.env.METRICS_PORT || "4042");
export const USER_PORT = parseInt(process.env.USER_PORT || "4043");
export const PACKAGE_PORT = parseInt(process.env.PACKAGE_PORT || "4044");
export const RATE_PORT = parseInt(process.env.RATE_PORT || "4045");
export const ORG_PORT = parseInt(process.env.ORG_PORT || "4046");
export const NETWORK_PORT = parseInt(process.env.NETWORK_PORT || "4047");
export const SUBSCRIBER_PORT = parseInt(process.env.SUBSCRIBER_PORT || "4048");
export const ALERT_PORT = parseInt(process.env.ALERT_PORT || "4049");
