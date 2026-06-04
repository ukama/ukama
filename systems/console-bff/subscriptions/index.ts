import { createYoga } from "graphql-yoga";
import { createServer } from "node:http";
import "reflect-metadata";
import { buildSchema } from "type-graphql";

import { SUBSCRIPTIONS_PORT } from "../common/configs";
import { logger } from "../common/logger";
import resolvers from "./resolvers";
import { pubSub } from "./resolvers/pubsub";

async function bootstrap() {
  const schema = await buildSchema({
    resolvers: resolvers,
    pubSub,
  });

  const yoga = createYoga({ schema });

  // Health endpoints for k8s probes and the api server's /ping check;
  // everything else goes to yoga (/graphql).
  const server = createServer((req, res) => {
    if (
      req.url === "/ping" ||
      req.url === "/healthz" ||
      req.url === "/readyz"
    ) {
      res.writeHead(200, { "content-type": "application/json" });
      res.end(JSON.stringify({ status: "ok" }));
      return;
    }
    yoga(req, res);
  });

  server.listen(SUBSCRIPTIONS_PORT, () => {
    logger.info(
      `Server is running on http://localhost:${SUBSCRIPTIONS_PORT}/graphql`
    );
  });
}
bootstrap().catch(logger.error);
