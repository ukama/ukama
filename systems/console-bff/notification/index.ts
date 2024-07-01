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
import NotificationAPI from "./datasource/notification_api";
import resolvers from "./resolvers";

const runServer = async () => {
  const isSuccess = await findProcessNKill(`${SUB_GRAPHS.notification.port}`);
  if (isSuccess) {
    const server = await SubGraphServer(resolvers);
    await startStandaloneServer(server, {
      context: async ({ req }) => {
        return {
          headers: parseGatewayHeaders(req.headers),
          dataSources: {
            dataSource: new NotificationAPI(),
          },
        };
      },
      listen: { port: SUB_GRAPHS.notification.port },
    });

    logger.info(
      `🚀 Ukama ${SUB_GRAPHS.notification.name} service running at http://localhost:${SUB_GRAPHS.notification.port}/graphql`
    );
  } else {
    logger.error(
      `Server failed to start on port ${SUB_GRAPHS.notification.port}`
    );
  }
};

runServer();
