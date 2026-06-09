/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import axios, { AxiosInstance, AxiosResponse } from "axios";

import { HTTP_TIMEOUT_MS } from "../configs";
import { ApiMethodDataDto } from "../types";

class ApiMethods {
  private readonly client: AxiosInstance;

  constructor() {
    this.client = axios.create({
      timeout: HTTP_TIMEOUT_MS,
    });
  }

  fetch = async (req: ApiMethodDataDto): Promise<AxiosResponse> => {
    return this.client.request({
      method: req.method,
      url: req.url,
      data: req.body,
      params: req.params,
      headers: req.headers,
      httpsAgent: req.httpsAgent,
    });
  };
}

export default new ApiMethods();
