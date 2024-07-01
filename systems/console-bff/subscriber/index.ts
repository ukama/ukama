/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { startStandaloneServer } from "@apollo/server/standalone";
import "reflect-metadata";

import { findProcessNKill, parseGatewayHeaders } from "../common/utils";
import SubGraphServer from "./../common/apollo";
import { SUB_GRAPHS } from "./../common/configs";
import { logger } from "./../common/logger";
import SubscriberAPI from "./datasource/subscriber_api";
import resolvers from "./resolver";

const runServer = async () => {
  const isSuccess = await findProcessNKill(`${SUB_GRAPHS.subscriber.port}`);
  if (isSuccess) {
    const server = await SubGraphServer(resolvers);
    await startStandaloneServer(server, {
      context: async ({ req }) => {
        return {
          headers: parseGatewayHeaders(req.headers),
          dataSources: {
            dataSource: new SubscriberAPI(),
          },
        };
      },
      listen: { port: SUB_GRAPHS.subscriber.port },
    });

    logger.info(
      `🚀 Ukama ${SUB_GRAPHS.subscriber.name} service running at http://localhost:${SUB_GRAPHS.subscriber.port}/graphql`
    );
  } else {
    logger.error(
      `Server failed to start on port ${SUB_GRAPHS.subscriber.port}`
    );
  }
};

runServer();
