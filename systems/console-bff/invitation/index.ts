import { startStandaloneServer } from "@apollo/server/standalone";
import "reflect-metadata";

import { parseGatewayHeaders } from "../common/utils";
import SubGraphServer from "./../common/apollo";
import { INVITATION_PORT } from "./../common/configs";
import { logger } from "./../common/logger";
import InvitationAPI from "./datasource/invitation_api";
import resolvers from "./resolver";

const runServer = async () => {
  const server = await SubGraphServer(resolvers);
  await startStandaloneServer(server, {
    context: async ({ req, res }) => {
      const { cache } = server;
      return {
        headers: parseGatewayHeaders(req.headers),
        dataSources: {
          dataSource: new InvitationAPI(),
        },
      };
    },
    listen: { port: INVITATION_PORT },
  });

  logger.info(
    `🚀 Ukama Invitation service running at http://localhost:${INVITATION_PORT}/graphql`
  );
};

runServer();
