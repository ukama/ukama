/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { startStandaloneServer } from "@apollo/server/standalone";
import "reflect-metadata";

import { openStore } from "../common/storage";
import { THeaders } from "../common/types";
import { getBaseURL, parseGatewayHeaders } from "../common/utils";
import SubGraphServer from "./../common/apollo";
import { SUB_GRAPHS } from "./../common/configs";
import { logger } from "./../common/logger";
import ReportAPI from "./datasource/report_api";
import resolvers from "./resolvers";

const runServer = async () => {
  const server = await SubGraphServer(resolvers);
  const store = openStore();
  await startStandaloneServer(server, {
    context: async ({ req }) => {
      const headers: THeaders = parseGatewayHeaders(req.headers);
      const baseURL = await getBaseURL(
        SUB_GRAPHS.report.name,
        headers.orgName,
        store
      );
      return {
        headers: headers,
        baseURL: baseURL.message,
        dataSources: {
          dataSource: new ReportAPI(),
        },
      };
    },

    listen: { port: SUB_GRAPHS.report.port },
  });

  logger.info(
    `🚀 Ukama ${SUB_GRAPHS.report.name} service running at http://localhost:${SUB_GRAPHS.report.port}/graphql`
  );
};

runServer();
