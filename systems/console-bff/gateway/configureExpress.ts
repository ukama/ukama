/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import express from "express";
import expressWinston from "express-winston";

import { rateLimit, securityHeaders } from "../common/middleware/security";

function configureExpress(logger: any) {
  const app = express();
  // Trust the proxy in front of the gateway so req.ip reflects the client.
  app.set("trust proxy", 1);
  app.disable("x-powered-by");

  app.use(securityHeaders());
  app.use(rateLimit());
  app.use(expressWinston.logger({ winstonInstance: logger }));
  return app;
}

export { configureExpress };
