/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Apps, HealthInfo } from "../resolvers/types";

/* The node app-health endpoint emits camelCase JSON at the REST boundary. */
interface AppResourceRest {
  cpuPercent: number;
  memoryRssKb: number;
  diskReadBytes: number;
  diskWriteBytes: number;
}
interface AppRest {
  name: string;
  version: string;
  tag: string;
  status: string;
  resource?: AppResourceRest | null;
}

export const mapApps = (res: { apps?: AppRest[] | null }): Apps => ({
  apps: (res.apps ?? []).map(app => ({
    name: app.name,
    version: app.version,
    tag: app.tag,
    status: app.status,
    resource: app.resource
      ? {
          cpuPercent: app.resource.cpuPercent,
          memoryRssKb: app.resource.memoryRssKb,
          diskReadBytes: app.resource.diskReadBytes,
          diskWriteBytes: app.resource.diskWriteBytes,
        }
      : undefined,
  })),
});

/* GET /v1/health/nodes/{nodeId}/reports returns the stored reports newest
 * first; payload is the node's raw health JSON, base64-encoded ([]byte). */
interface HealthReportRest {
  id?: string;
  nodeId?: string;
  reportedAt?: number;
  payload?: string;
}

interface HealthPayload {
  system?: { uptimeSec?: number; starter?: { state?: string } };
  interfaces?: {
    radio?: { state?: string } | null;
    cellular?: { service?: string } | null;
  };
  apps?: {
    space?: string;
    name?: string;
    tag?: string;
    state?: string;
    resources?: Record<string, number>;
  }[];
}

const decodePayload = (payload?: string): HealthPayload => {
  if (!payload) return {};
  try {
    return JSON.parse(Buffer.from(payload, "base64").toString("utf-8"));
  } catch {
    return {};
  }
};

export const mapReportsToHealthInfo = (res: {
  reports?: HealthReportRest[] | null;
}): HealthInfo => {
  const report = res?.reports?.[0] ?? {};
  const reportId = report.id ?? "";
  const payload = decodePayload(report.payload);

  // Name/value pairs kept for consumers of the older system[] shape, which
  // read the "radio" and "service" entries.
  const systemValues: [string, string | number | undefined][] = [
    ["radio", payload.interfaces?.radio?.state],
    ["service", payload.interfaces?.cellular?.service],
    ["uptimeSec", payload.system?.uptimeSec],
    ["starter", payload.system?.starter?.state],
  ];
  const system = systemValues
    .filter(([, value]) => value !== undefined && value !== "")
    .map(([name, value]) => ({
      id: `${reportId}-${name}`,
      healthId: reportId,
      name,
      value: String(value),
    }));

  const capps = (payload.apps ?? []).map(app => {
    const cappId = `${reportId}-${app.name ?? ""}`;
    return {
      id: cappId,
      space: app.space ?? "",
      name: app.name ?? "",
      tag: app.tag ?? "",
      status: app.state ?? "",
      resources: Object.entries(app.resources ?? {}).map(([name, value]) => ({
        id: `${cappId}-${name}`,
        cappId,
        name,
        value: String(value),
      })),
    };
  });

  return {
    id: reportId,
    nodeId: report.nodeId ?? "",
    timestamp: report.reportedAt ? String(report.reportedAt) : "",
    system,
    capps,
  };
};
