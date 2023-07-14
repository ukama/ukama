import { ApolloGateway, IntrospectAndCompose } from "@apollo/gateway";
import { ApolloServer } from "@apollo/server";
import { expressMiddleware } from "@apollo/server/express4";
import { ApolloServerPluginDrainHttpServer } from "@apollo/server/plugin/drainHttpServer";
import { ApolloServerPluginInlineTrace } from "@apollo/server/plugin/inlineTrace";
import { json } from "body-parser";
import cors from "cors";
import { createServer } from "http";

import {
  AUTH_APP_URL,
  CONSOLE_APP_URL,
  GATEWAY_PORT,
  PLANNING_SERVICE_PORT,
  PLAYGROUND_URL,
} from "../common/configs";
import { logger } from "../common/logger";
import { configureExpress } from "./configureExpress";

const app = configureExpress(logger);
const httpServer = createServer(app);

const startServer = async () => {
  const gateway = new ApolloGateway({
    supergraphSdl: new IntrospectAndCompose({
      subgraphs: [
        // { name: "metrics", url: `http://localhost:${METRICS_PORT}` },
        { name: "planning", url: `http://localhost:${PLANNING_SERVICE_PORT}` },
      ],
    }),
  });
  await gateway.load();
  const server = new ApolloServer({
    gateway,
    plugins: [
      ApolloServerPluginInlineTrace({}),
      ApolloServerPluginDrainHttpServer({ httpServer }),
    ],
  });

  await server.start();

  app.use(
    "/graphql",
    cors({
      origin: [AUTH_APP_URL, PLAYGROUND_URL, CONSOLE_APP_URL],
      credentials: true,
    }),

    json(),
    expressMiddleware(server)
  );
  await new Promise((resolve: any) =>
    httpServer.listen({ port: GATEWAY_PORT }, resolve)
  );
  logger.info(`🚀 Server ready at http://localhost:${GATEWAY_PORT}/graphql`);
};

startServer();
