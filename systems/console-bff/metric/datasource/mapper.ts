/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import {
  GetNodeLatestMetricInput,
  GetSiteLatestMetricInput,
  NodeLatestMetric,
  SiteLatestMetric,
} from "../resolver/types";

const parseLatestMetricRes = (
  res: any,
  args: { type: string; [key: string]: any },
  entityType: "node" | "site"
) => {
  const data = res.data.result[0];
  const idKey = `${entityType}Id`;
  const metricResponse = {
    success: true,
    msg: "success",
  };

  if (data?.value?.length > 0) {
    return {
      ...metricResponse,
      orgId: data.metric.org,
      [idKey]: args[idKey],
      type: args.type,
      value: data.value,
    };
  } else {
    return {
      ...metricResponse,
      orgId: "",
      [idKey]: "",
      type: "",
      value: [0, 0],
    };
  }
};

export const parseNodeLatestMetricRes = (
  res: any,
  args: GetNodeLatestMetricInput
): NodeLatestMetric => {
  return parseLatestMetricRes(res, args, "node") as NodeLatestMetric;
};

export const parseSiteLatestMetricRes = (
  res: any,
  args: GetSiteLatestMetricInput
): SiteLatestMetric => {
  return parseLatestMetricRes(res, args, "site") as SiteLatestMetric;
};
