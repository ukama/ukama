/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Shared express bootstrap: request-id, ALS request context, security
 * headers, rate limiting and structured access logs. Used by the gateway
 * today and the consolidated API server (server/) going forward.
 */
import express, { type Request } from "express";
import expressWinston from "express-winston";

import { runWithRequestId } from "../logger/requestContext";
import { compression } from "./compression";
import { REQUEST_ID_HEADER, requestId } from "./requestId";
import { rateLimit, securityHeaders } from "./security";

function configureExpress(logger: any) {
  const app = express();
  // Trust the proxy in front of the server so req.ip reflects the client.
  app.set("trust proxy", 1);
  app.disable("x-powered-by");

  app.use(requestId());
  // Bind the correlation id to the async context for the whole request so
  // every downstream log line (resolvers, datasources) carries it.
  app.use((req, _res, next) =>
    runWithRequestId(req.headers[REQUEST_ID_HEADER] as string, () => next())
  );
  app.use(securityHeaders());
  app.use(rateLimit());
  // Gzip buffered JSON responses (GraphQL payloads) above the size threshold.
  app.use(compression());
  app.use(
    expressWinston.logger({
      winstonInstance: logger,
      dynamicMeta: (req: Request) => ({
        requestId: req.headers[REQUEST_ID_HEADER],
      }),
    })
  );
  return app;
}

export { configureExpress };
