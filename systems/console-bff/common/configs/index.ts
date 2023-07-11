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
