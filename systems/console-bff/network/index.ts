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

import { THeaders } from "../common/types";
import {
  findProcessNKill,
  getBaseURL,
  parseGatewayHeaders,
} from "../common/utils";
import SubGraphServer from "./../common/apollo";
import { SUB_GRAPHS } from "./../common/configs";
import { logger } from "./../common/logger";
import NetworkAPI from "./datasource/network_api";
import resolvers from "./resolvers";

const runServer = async () => {
  const isSuccess = await findProcessNKill(`${SUB_GRAPHS.network.port}`);
  if (isSuccess) {
    const server = await SubGraphServer(resolvers);
    const redisClient = createClient().on("error", error => {
      logger.error(
        `Error creating redis for ${SUB_GRAPHS.network.name} service, Error: ${error}`
      );
    });
    const connectPromise = redisClient.connect();
    await connectPromise;

    await startStandaloneServer(server, {
      context: async ({ req }) => {
        const headers: THeaders = parseGatewayHeaders(req.headers);
        const baseURL = await getBaseURL(
          SUB_GRAPHS.member.name,
          headers.orgName,
          redisClient.isOpen ? redisClient : null
        );
        return {
          headers: headers,
          baseURL: baseURL.message,
          dataSources: {
            dataSource: new NetworkAPI(),
          },
        };
      },
      listen: { port: SUB_GRAPHS.network.port },
    });

    logger.info(
      `🚀 Ukama ${SUB_GRAPHS.network.name} service running at http://localhost:${SUB_GRAPHS.network.port}/graphql`
    );
  } else {
    logger.error(`Server failed to start on port ${SUB_GRAPHS.network.port}`);
  }
};

runServer();
