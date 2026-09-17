/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { BaseRESTDataSource } from "../../common/datasource";
import { GetNodeLatestMetricInput, NodeLatestMetric } from "../resolver/types";
import { parseNodeLatestMetricRes } from "./mapper";

const VERSION = "v1";
const METRICS = "metrics";

class MetricAPI extends BaseRESTDataSource {
  getNodeLatestMetric = async (
    baseURL: string,
    args: GetNodeLatestMetricInput
  ): Promise<NodeLatestMetric> => {
    this.logger.info(
      `GetNodeLatestMetric [GET]: ${baseURL}/${VERSION}/${METRICS}/${args.type}`
    );
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/${METRICS}/${args.type}`).then(res =>
      parseNodeLatestMetricRes(res, args)
    );
  };

  /**
   * Generic latest-value read for one metric key (org-scoped, no entity
   * stamping). Used by the dashboard KPI sections (plan Phase 4 — polled,
   * no subscriptions in v1).
   *
   * NOTE: this hits `/v1/metrics/:metric`, whose gateway handler hardcodes the
   * `system` node type. Only use it for system/org-scoped keys — for per-node
   * KPIs use getNodeLatest (the gateway can't resolve a node type here).
   */
  getLatestMetric = async (
    baseURL: string,
    type: string
  ): Promise<{ type: string; value: [number, number]; success: boolean }> => {
    this.logger.info(
      `GetLatestMetric [GET]: ${baseURL}/${VERSION}/${METRICS}/${type}`
    );
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/${METRICS}/${type}`).then(res => {
      const data = res?.data?.result?.[0];
      if (data?.value?.length > 0) {
        return { type, value: data.value as [number, number], success: true };
      }
      return { type, value: [0, 0] as [number, number], success: false };
    });
  };

  /** Lookback for the `/v1/last` read. The pushgateway keeps serving a dead
   *  node's last sample, so a caller must gate on connectivity, not on this. */
  private static readonly LAST_LOOKBACK = "1h";

  /**
   * Latest sample for a per-node metric via the gateway's instant endpoint
   * (`/v1/last/metrics/:metric?node=`), which resolves the node type from the
   * id and falls back to the `system` bucket for keys such as com_uptime and
   * ctl_uptime. One value per series; the first series is taken.
   */
  getNodeLast = async (
    baseURL: string,
    type: string,
    nodeId: string
  ): Promise<{ type: string; value: [number, number]; success: boolean }> => {
    const path = `/${VERSION}/last/${METRICS}/${type}?node=${nodeId}&fn=last&lookback=${MetricAPI.LAST_LOOKBACK}`;
    this.logger.info(`GetNodeLast [GET]: ${baseURL}${path}`);
    this.baseURL = baseURL;
    return this.get(path).then(raw => {
      const res = typeof raw === "string" ? JSON.parse(raw) : raw;
      const value = res?.data?.result?.[0]?.value as
        | [number, string]
        | undefined;
      if (Array.isArray(value) && value.length === 2) {
        return {
          type,
          value: [Number(value[0]), Number(value[1])] as [number, number],
          success: Number.isFinite(Number(value[1])),
        };
      }
      return { type, value: [0, 0] as [number, number], success: false };
    });
  };

  /** Lookback window (seconds) for deriving a node's latest KPI value. */
  private static readonly LATEST_LOOKBACK = 3600;
  private static readonly LATEST_STEP = 60;

  /**
   * Latest value for a per-node metric. The gateway has no node-scoped instant
   * endpoint, so we query the node range endpoint (`/v1/nodes/:node/metrics/:metric`,
   * which resolves node type from the id) over a short recent window and take
   * the most recent sample.
   */
  getNodeLatest = async (
    baseURL: string,
    type: string,
    nodeId: string
  ): Promise<{ type: string; value: [number, number]; success: boolean }> => {
    const to = Math.floor(Date.now() / 1000);
    const from = to - MetricAPI.LATEST_LOOKBACK;
    const path = `/${VERSION}/nodes/${nodeId}/${METRICS}/${type}?from=${from}&to=${to}&step=${MetricAPI.LATEST_STEP}`;
    this.logger.info(`GetNodeLatest [GET]: ${baseURL}${path}`);
    this.baseURL = baseURL;
    return this.get(path).then(res => {
      const values = res?.data?.result?.[0]?.values as
        | [number, string][]
        | undefined;
      const last = Array.isArray(values) ? values[values.length - 1] : null;
      if (last && last.length === 2) {
        return {
          type,
          value: [Number(last[0]), Number(last[1])] as [number, number],
          success: true,
        };
      }
      return { type, value: [0, 0] as [number, number], success: false };
    });
  };
}

export default MetricAPI;
