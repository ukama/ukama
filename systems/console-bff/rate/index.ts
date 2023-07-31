import { startStandaloneServer } from "@apollo/server/standalone";
import "reflect-metadata";

import SubGraphServer from "./../common/apollo";
import { RATE_PORT } from "./../common/configs";
import { logger } from "./../common/logger";
import RateAPI from "./datasource/rate_api";
import resolvers from "./resolver";

const runServer = async () => {
  const server = await SubGraphServer(resolvers);
  await startStandaloneServer(server, {
    context: async () => {
      const { cache } = server;
      return {
        dataSources: {
          dataSource: new RateAPI(),
        },
      };
    },
    listen: { port: RATE_PORT },
  });

  logger.info(
    `🚀 Ukama Node service running at http://localhost:${RATE_PORT}/graphql`
  );
};

runServer();
