/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { GraphQLError } from "graphql";

/**
 * Error for failed authentication. Carries the GraphQL `UNAUTHENTICATED`
 * code and an HTTP 401 status so Apollo Server returns a proper 401
 * response (instead of a generic 500) when thrown from a context
 * function or resolver. Clients rely on the 401 to trigger re-login.
 */
export class UnauthenticatedError extends GraphQLError {
  constructor(message: string) {
    super(message, {
      extensions: {
        code: "UNAUTHENTICATED",
        http: { status: 401 },
      },
    });
  }
}
