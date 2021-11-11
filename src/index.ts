import "reflect-metadata";
import setupLogger from "./config/logger";
import configureExpress from "./config/express";
import configureApolloServer from "./config/apolloServer";
import env from "./config/env";
import { mockServer } from "./mockServer";
const logger = setupLogger("app");
const { PORT } = env || 3000;

const initializeApp = async () => {
    const app = configureExpress({
        logger,
    });

    const server = await configureApolloServer();
    server.applyMiddleware({ app });

    mockServer(app);
    app.listen(PORT, () => logger.info(`Server listening on port: ${PORT}`));
};

initializeApp().catch(error => logger.error(error));
