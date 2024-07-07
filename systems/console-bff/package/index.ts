/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { startStandaloneServer } from "@apollo/server/standalone";
import { createClient } from "redis";
import "reflect-metadata";

import SubGraphServer from "../common/apollo";
import { BFF_REDIS, SUB_GRAPHS } from "../common/configs";
import { logger } from "../common/logger";
import { THeaders } from "../common/types";
import { getBaseURL, parseGatewayHeaders } from "../common/utils";
import PackageAPI from "./datasource/package_api";
import resolvers from "./resolver";

const runServer = async () => {
  const server = await SubGraphServer(resolvers);
  const redisClient = createClient({
    url: BFF_REDIS,
  }).on("error", error => {
    logger.error(
      `Error creating redis for ${SUB_GRAPHS.package.name} service, Error: ${error}`
    );
  });
  const connectPromise = redisClient.connect();
  await connectPromise;

  await startStandaloneServer(server, {
    context: async ({ req }) => {
      const headers: THeaders = parseGatewayHeaders(req.headers);
      const baseURL = await getBaseURL(
        SUB_GRAPHS.package.name,
        headers.orgName,
        redisClient.isOpen ? redisClient : null
      );
      return {
        headers: headers,
        baseURL: baseURL.message,
        dataSources: {
          dataSource: new PackageAPI(),
        },
      };
    },
    listen: { port: SUB_GRAPHS.package.port },
  });

  logger.info(
    `🚀 Ukama ${SUB_GRAPHS.package.name} service running at http://localhost:${SUB_GRAPHS.package.port}/graphql`
  );
};

runServer();
