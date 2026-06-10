/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { THeaders } from "../../common/types";
import AnalyticsAPI from "../datasource/analytics_api";

export interface Context {
  baseURL: string;
  dataSources: {
    analytics: AnalyticsAPI;
  };
  headers: THeaders;
}
