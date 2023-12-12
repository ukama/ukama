/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import axios from "axios";

import { ApiMethodDataDto } from "../types";

class ApiMethods {
  constructor() {
    axios.create({
      timeout: 10000,
    });
  }
  fetch = async (req: ApiMethodDataDto) => {
    return axios(req as any).catch(err => {
      throw err;
    });
  };
}

export default new ApiMethods();
