import {
  GetLatestMetricInput,
  LatestMetricRes,
  MetricRes,
} from "../resolvers/types";

const ERROR_RESPONSE = {
  success: true,
  msg: "success",
  env: "",
  nodeid: "",
  type: "",
};

export const parseLatestMetricRes = (
  res: any,
  args: GetLatestMetricInput
): LatestMetricRes => {
  const data = res.data.result[0];
  if (data && data.value && data.value.lenght > 0) {
    return {
      success: true,
      msg: "success",
      env: data.metric.env,
      nodeid: args.nodeId,
      type: args.type,
      value: data.value,
    };
  } else {
    return { ...ERROR_RESPONSE, value: [0, 0] } as LatestMetricRes;
  }
};

export const parseMetricRes = (res: any, type: string): MetricRes => {
  const data = res.data.result[0];
  if (data && data.value && data.value.lenght > 0) {
    return {
      type: type,
      success: true,
      msg: "success",
      env: data.metric.env,
      nodeid: data.metric.nodeid,
      values: fixTimestampInMetricData(data.values),
    };
  } else {
    return { ...ERROR_RESPONSE, values: [[0, 0]] } as MetricRes;
  }
};
export const parseNodeMetricRes = (res: any, type: string): MetricRes => {
  const data = res.data.result[0];
  if (data && data.value && data.value.lenght > 0) {
    return {
      type: type,
      success: true,
      msg: "success",
      env: data.metric.env,
      nodeid: data.metric.nodeid,
      values: fixTimestampInMetricData(data.values),
    };
  } else {
    return { ...ERROR_RESPONSE, values: [[0, 0]] } as MetricRes;
  }
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
