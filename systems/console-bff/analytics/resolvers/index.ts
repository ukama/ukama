/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { NonEmptyArray } from "type-graphql";

import { GetKpiValuesResolver } from "./getKpiValues";
import { GetPerformanceReportResolver } from "./getPerformanceReport";

/**
 * Analytics subgraph — a thin GraphQL surface over the analytics gateway's
 * generic KPI read API. The gateway exposes only self-describing KPI endpoints
 * (see systems/analytics/api-gateway): latest KPI values and resource
 * performance reports. Add resolvers here as more of that API is surfaced.
 */
const resolvers: NonEmptyArray<any> = [
  GetKpiValuesResolver,
  GetPerformanceReportResolver,
];

export default resolvers;
