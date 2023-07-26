import {
  GetLatestMetricInput,
  LatestMetricRes,
  MetricRes,
} from "../resolvers/types";

export const parseLatestMetricRes = (
  res: any,
  args: GetLatestMetricInput
): LatestMetricRes => {
  const data = res.data.result[0];
  return {
    env: data.metric.env,
    nodeid: args.nodeId,
    type: args.type,
    value: data.value,
  };
};
export const parseMetricRes = (
  res: any,
  args: GetLatestMetricInput
): MetricRes => {
  const data = res.data.result[0];
  return {
    env: data.metric.env,
    nodeid: args.nodeId,
    type: args.type,
    values: data.values,
  };
};
