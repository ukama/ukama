import { startStandaloneServer } from "@apollo/server/standalone";
import "reflect-metadata";

import { parseHeaders } from "../common/utils";
import SubGraphServer from "./../common/apollo";
import { PACKAGE_PORT } from "./../common/configs";
import { logger } from "./../common/logger";
import PackageAPI from "./datasource/package_api";
import resolvers from "./resolver";

const runServer = async () => {
  const server = await SubGraphServer(resolvers);
  await startStandaloneServer(server, {
    context: async ({ req, res }) => {
      const { cache } = server;
      return {
        headers: parseHeaders(req.headers),
        dataSources: {
          dataSource: new PackageAPI(),
        },
      };
    },
    listen: { port: PACKAGE_PORT },
  });

  logger.info(
    `🚀 Ukama Node service running at http://localhost:${PACKAGE_PORT}/graphql`
  );
};

runServer();
