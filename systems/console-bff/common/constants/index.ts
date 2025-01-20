/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { getGraphsKeyByType } from "../utils";
import { GRAPHS_TYPE } from "./../enums/index";

const getAllMetricsKeys = (nodeId: string): string[] => {
  const allKeys: string[] = [];
  const types = [
    GRAPHS_TYPE.NODE_HEALTH,
    GRAPHS_TYPE.NETWORK,
    GRAPHS_TYPE.RESOURCES,
    GRAPHS_TYPE.RADIO,
    GRAPHS_TYPE.SUBSCRIBERS,
  ];

  types.forEach(type => {
    const keys = getGraphsKeyByType(type, nodeId);
    keys.forEach(key => {
      if (!allKeys.includes(key)) {
        allKeys.push(key);
      }
    });
  });

  return allKeys;
};

export { getAllMetricsKeys };
