/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { HealthInfo } from "../resolvers/types";

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
