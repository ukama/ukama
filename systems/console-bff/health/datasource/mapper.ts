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

export const dtoToHealthInfo = (res: any): HealthInfo => {
  const health = res?.healths?.[0] ?? {};

  const system = (health.system ?? []).map((item: any) => ({
    id: item.id,
    healthId: item.healthId,
    name: item.name,
    value: item.value,
  }));

  const capps = (health.capps ?? []).map((item: any) => ({
    id: item.id,
    space: item.space,
    name: item.name,
    tag: item.tag,
    status: item.status,
    resources: (item.resources ?? []).map((resource: any) => ({
      id: resource.id,
      cappId: resource.cappId,
      name: resource.name,
      value: resource.value,
    })),
  }));

  return {
    id: health.id ?? "",
    nodeId: health.nodeId ?? "",
    timestamp: health.timestamp ?? "",
    system,
    capps,
  };
};
