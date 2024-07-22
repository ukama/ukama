/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { GetNodeLatestMetricInput, NodeLatestMetric } from "../resolver/types";

const ERROR_RESPONSE = {
  success: true,
  msg: "success",
  orgId: "",
  nodeId: "",
  type: "",
};

export const parseNodeLatestMetricRes = (
  res: any,
  args: GetNodeLatestMetricInput
): NodeLatestMetric => {
  const data = res.data.result[0];
  if (data?.value?.length > 0) {
    return {
      success: true,
      msg: "success",
      orgId: data.metric.org,
      nodeId: args.nodeId,
      type: args.type,
      value: data.value,
    };
  } else {
    return { ...ERROR_RESPONSE, value: [0, 0] } as NodeLatestMetric;
  }
};
