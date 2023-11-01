/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { startStandaloneServer } from "@apollo/server/standalone";
import "reflect-metadata";

import SubGraphServer from "../common/apollo";
import { PACKAGE_PORT } from "../common/configs";
import { logger } from "../common/logger";
import { parseGatewayHeaders } from "../common/utils";
import PackageAPI from "./datasource/package_api";
import resolvers from "./resolver";

const runServer = async () => {
  const server = await SubGraphServer(resolvers);
  await startStandaloneServer(server, {
    context: async ({ req }) => {
      return {
        headers: parseGatewayHeaders(req.headers),
        dataSources: {
          dataSource: new PackageAPI(),
        },
      };
    },
    listen: { port: PACKAGE_PORT },
  });

  logger.info(
    `🚀 Ukama User service running at http://localhost:${PACKAGE_PORT}/graphql`
  );
};

runServer();
