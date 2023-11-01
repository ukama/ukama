/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { startStandaloneServer } from "@apollo/server/standalone";
import "reflect-metadata";

import { parseHeaders } from "../common/utils";
import SubGraphServer from "./../common/apollo";
import { ALERT_PORT } from "./../common/configs";
import { logger } from "./../common/logger";
import AlertAPI from "./datasource/alert_api";
import resolvers from "./resolver";

const runServer = async () => {
  const server = await SubGraphServer(resolvers);
  await startStandaloneServer(server, {
    context: async ({ req }) => {
      // const { cache } = server;
      return {
        headers: parseHeaders(req.headers),
        dataSources: {
          dataSource: new AlertAPI(),
        },
      };
    },
    listen: { port: ALERT_PORT },
  });

  logger.info(
    `🚀 Ukama Node service running at http://localhost:${ALERT_PORT}/graphql`
  );
};

runServer();
