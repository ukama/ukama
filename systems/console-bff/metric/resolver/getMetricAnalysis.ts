/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Query, Resolver } from "type-graphql";
import { logger } from "../../common/logger";
import { MetricAnalysis } from "./types";

@Resolver()
export class GetMetricAnalysisResolver {
  @Query(() => MetricAnalysis)
  async getMetricAnalysis(
    @Arg("metricId") metricId: string,
    @Arg("nodeId") nodeId: string
  ): Promise<MetricAnalysis> {
    logger.info(
      `Getting metric analysis for metric ${metricId} and node ${nodeId}`
    );
    return {
      aggregated: {
        computed_at: new Date().toISOString(),
        value: 0,
        min: 0,
        max: 0,
        p95: 0,
        mean: 0,
        median: 0,
        sample_count: 0,
        aggregation: "average",
        noise_estimate: 0,
      },
      trend: {
        type: "average",
        value: 0,
      },
      confidence: {
        value: 0,
      },
      projection: {
        type: "average",
        eta_sec: 0,
      },
      state: "ok",
    };
  }
}
