import { createYoga } from "graphql-yoga";
import { createServer } from "node:http";
import "reflect-metadata";
import { buildSchema } from "type-graphql";

import { SUBSCRIPTIONS_PORT } from "../common/configs";
import { logger } from "../common/logger";
import { openStore } from "../common/storage";
import resolvers from "./resolvers";
import { pubSub } from "./resolvers/pubsub";

async function bootstrap() {
  const schema = await buildSchema({
    resolvers: resolvers,
    pubSub,
  });

  const store = openStore();
  const yoga = createYoga({
    schema,
    context: {
      store: store,
    },
  });

  const server = createServer(yoga);

  server.listen(SUBSCRIPTIONS_PORT, () => {
    logger.info(
      `Server is running on http://localhost:${SUBSCRIPTIONS_PORT}/graphql`
    );
  });
}
bootstrap().catch(logger.error);
