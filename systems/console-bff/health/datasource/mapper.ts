/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { HealthInfo } from "../resolvers/types";

export const dtoToHealthInfo = (res: HealthInfo): HealthInfo => {
  return {
    id: res.id,
    nodeId: res.nodeId,
    timestamp: res.timestamp,
    system: res.system,
    capps: res.capps,
  };
};
