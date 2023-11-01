/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { BaseError } from "./base.error";
import { HttpStatusCode } from "./codes";

export class HTTP500Error extends BaseError {
  constructor(description: string) {
    super("INTERNAL SERVER", HttpStatusCode.INTERNAL_SERVER, description, true);
  }
}
