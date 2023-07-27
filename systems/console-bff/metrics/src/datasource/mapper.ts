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
    values: fixTimestampInMetricData(data.values),
  };
};
export const parseNodeMetricRes = (res: any, type: string): MetricRes => {
  const data = res.data.result[0];
  return {
    type: type,
    values: fixTimestampInMetricData(data.values),
    env: data.metric.env,
    nodeid: data.metric.nodeid,
  };
};

const fixTimestampInMetricData = (
  values: [[number, string]]
): [number, number][] => {
  if (values.length > 0) {
    const fixedValues: [number, number][] = values.map(
      (value: [number, string]) => {
        return [Math.floor(value[0]) * 1000, parseFloat(value[1])];
      }
    );
    return fixedValues;
  }
  return [];
};
