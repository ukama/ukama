import { startStandaloneServer } from "@apollo/server/standalone";
import "reflect-metadata";

import SubGraphServer from "./../common/apollo";
import { BILLING_PORT } from "./../common/configs";
import { logger } from "./../common/logger";
import BillingAPI from "./datasource/billing_api";
import resolvers from "./resolvers";

const runServer = async () => {
  const server = await SubGraphServer(resolvers);
  await startStandaloneServer(server, {
    context: async () => {
      const { cache } = server;
      return {
        dataSources: {
          dataSource: new BillingAPI(),
        },
      };
    },
    listen: { port: BILLING_PORT },
  });

  logger.info(
    `🚀 Ukama Node service running at http://localhost:${BILLING_PORT}/graphql`
  );
};

runServer();
