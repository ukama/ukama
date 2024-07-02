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
import { COMPONENT_INVENTORY_PORT } from "../common/configs";
import { logger } from "../common/logger";
import { parseGatewayHeaders } from "../common/utils";
import ComponentAPI from "./datasource/component_api";
import resolvers from "./resolvers";

const runServer = async () => {
  const server = await SubGraphServer(resolvers);
  await startStandaloneServer(server, {
    context: async ({ req }) => {
      return {
        headers: parseGatewayHeaders(req.headers),
        dataSources: {
          dataSource: new ComponentAPI(),
        },
      };
    },
    listen: { port: COMPONENT_INVENTORY_PORT },
  });

  logger.info(
    `🚀 Ukama Component Inventory service running at http://localhost:${COMPONENT_INVENTORY_PORT}/graphql`
  );
};

runServer();
