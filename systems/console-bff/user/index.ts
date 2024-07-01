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
import UserAPI from "./datasource/user_api";
import resolvers from "./resolver";

const runServer = async () => {
  const isSuccess = await findProcessNKill(`${SUB_GRAPHS.user.port}`);
  if (isSuccess) {
    const server = await SubGraphServer(resolvers);
    const redisClient = createClient();
    const connectPromise = redisClient.connect();
    await connectPromise;

    await startStandaloneServer(server, {
      context: async ({ req }) => {
        const headers: THeaders = parseGatewayHeaders(req.headers);
        const baseURL = getBaseURL(
          SUB_GRAPHS.user.name,
          headers.orgName,
          redisClient.isOpen ? redisClient : null
        );
        return {
          headers: headers,
          baseURL: baseURL,
          dataSources: {
            dataSource: new UserAPI(),
          },
        };
      },
      listen: { port: SUB_GRAPHS.user.port },
    });

    logger.info(
      `🚀 Ukama ${SUB_GRAPHS.user.name} service running at http://localhost:${SUB_GRAPHS.user.port}/graphql`
    );
  } else {
    logger.error(`Server failed to start on port ${SUB_GRAPHS.user.port}`);
  }
};

runServer();
