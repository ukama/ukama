import { GraphQLError } from "graphql";

import { asyncRestCall } from "../../common/axiosClient";
import { METRIC_API_GW } from "../../common/configs";
import { API_METHOD_TYPE } from "../../common/enums";
import {
  GetLatestMetricInput,
  GetMetricRangeInput,
  LatestMetricRes,
  MetricRes,
} from "../resolvers/types";
import {
  parseLatestMetricRes,
  parseMetricRes,
  parseNodeMetricRes,
} from "./mapper";

const getLatestMetric = async (
  args: GetLatestMetricInput
): Promise<LatestMetricRes> => {
  return await asyncRestCall({
    method: API_METHOD_TYPE.GET,
    url: `${METRIC_API_GW}/v1/metrics/${args.type}`,
  }).then(res => parseLatestMetricRes(res.data, args));
};

const getMetricRange = async (
  args: GetMetricRangeInput
): Promise<MetricRes> => {
  const { from, to = 0, step = 1 } = args;
  return await asyncRestCall({
    method: API_METHOD_TYPE.GET,
    url: `${METRIC_API_GW}/v1/range/metrics/${args.type}?from=${from}&to=${to}&step=${step}`,
  })
    .then(res => parseMetricRes(res.data, args.type))
    .catch(err => {
      throw new GraphQLError(err);
    });
};

const getNodeRangeMetric = async (
  args: GetMetricRangeInput
): Promise<MetricRes> => {
  const { from, to = 0, step = 1 } = args;
  console.log(
    "URL:",
    `${METRIC_API_GW}/v1/nodes/${args.nodeId}/metrics/${args.type}?from=${from}&to=${to}&step=${step}`
  );
  return await asyncRestCall({
    method: API_METHOD_TYPE.GET,
    url: `${METRIC_API_GW}/v1/nodes/${args.nodeId}/metrics/${args.type}?from=${from}&to=${to}&step=${step}`,
  })
    .then(res => {
      console.log("RES: ", res);
      return parseNodeMetricRes(res.data, args.type);
    })
    .catch(err => {
      console.log("ERR:", err);
      throw new GraphQLError(err);
    });
};

export { getLatestMetric, getMetricRange, getNodeRangeMetric };
