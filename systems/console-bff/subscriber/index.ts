import { startStandaloneServer } from "@apollo/server/standalone";
import "reflect-metadata";

import { parseHeaders } from "../common/utils";
import SubGraphServer from "./../common/apollo";
import { SUBSCRIBER_PORT } from "./../common/configs";
import { logger } from "./../common/logger";
import SubscriberAPI from "./datasource/subscriber_api";
import resolvers from "./resolver";

const runServer = async () => {
  const server = await SubGraphServer(resolvers);
  await startStandaloneServer(server, {
    context: async ({ req, res }) => {
      const { cache } = server;
      return {
        headers: parseHeaders(req.headers),
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
