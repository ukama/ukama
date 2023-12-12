/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { MiddlewareInterface, NextFn, ResolverData } from "type-graphql";
import { Service } from "typedi";

import { HTTP401Error, Messages } from "../errors";

@Service()
export class Authentication implements MiddlewareInterface<any> {
  async use({ context }: ResolverData<any>, next: NextFn): Promise<void> {
    if (context.req.headers !== undefined) {
      const token = context.req.headers["x-session-token"] ?? "";
      const cookie =
        context.req.headers.cookie &&
        context.req.headers.cookie.includes("ukama_session")
          ? context.req.headers.cookie
          : "";

      if (!cookie && !token) {
        throw new HTTP401Error(Messages.REQUEST_AUTHENTICATION_FAILED);
      }
      context.authType = token ? "token" : "cookie";
    }
    return next();
  }
}
