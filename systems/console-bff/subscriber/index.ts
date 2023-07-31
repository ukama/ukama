import { startStandaloneServer } from "@apollo/server/standalone";
import "reflect-metadata";

import SubGraphServer from "./../common/apollo";
import { SUBSCRIBER_PORT } from "./../common/configs";
import { logger } from "./../common/logger";
import SubscriberAPI from "./datasource/subscriber_api";
import resolvers from "./resolver";

const runServer = async () => {
  const server = await SubGraphServer(resolvers);
  await startStandaloneServer(server, {
    context: async () => {
      const { cache } = server;
      return {
        dataSources: {
          dataSource: new SubscriberAPI(),
        },
      };
    },
    listen: { port: SUBSCRIBER_PORT },
  });

  logger.info(
    `🚀 Ukama Node service running at http://localhost:${SUBSCRIBER_PORT}/graphql`
  );
};

runServer();
