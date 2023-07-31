import { startStandaloneServer } from "@apollo/server/standalone";
import "reflect-metadata";

import SubGraphServer from "./../common/apollo";
import { ALERT_PORT } from "./../common/configs";
import { logger } from "./../common/logger";
import AlertAPI from "./datasource/alert_api";
import resolvers from "./resolver";

const runServer = async () => {
  const server = await SubGraphServer(resolvers);
  await startStandaloneServer(server, {
    context: async () => {
      const { cache } = server;
      return {
        dataSources: {
          dataSource: new AlertAPI(),
        },
      };
    },
    listen: { port: ALERT_PORT },
  });

  logger.info(
    `🚀 Ukama Node service running at http://localhost:${ALERT_PORT}/graphql`
  );
};

runServer();
