/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { RESTDataSource } from "@apollo/datasource-rest";
import type {
  Fetcher,
  FetcherRequestInit,
  FetcherResponse,
} from "@apollo/utils.fetcher";

import { HTTP_TIMEOUT_MS } from "../configs";

/**
 * A fetch wrapper that aborts any upstream HTTP call exceeding
 * HTTP_TIMEOUT_MS, so a slow backend can never hang the BFF.
 */
const fetchWithTimeout: Fetcher = async (
  url: string,
  init?: FetcherRequestInit
): Promise<FetcherResponse> => {
  const controller = new AbortController();
  const timer = setTimeout(() => controller.abort(), HTTP_TIMEOUT_MS);
  try {
    return await fetch(url, {
      ...(init ?? {}),
      signal: controller.signal,
    } as RequestInit);
  } finally {
    clearTimeout(timer);
  }
};

/**
 * Base class for all subgraph datasources. Extends Apollo's
 * RESTDataSource with a hard request timeout. All datasources
 * must extend this class instead of RESTDataSource directly.
 */
export class BaseRESTDataSource extends RESTDataSource {
  constructor() {
    super({ fetch: fetchWithTimeout });
  }
}
