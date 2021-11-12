import "reflect-metadata";
import setupLogger from "./config/logger";
import configureExpress from "./config/express";
import configureApolloServer from "./config/apolloServer";
import { mockServer } from "./mockServer";
import { PORT } from "./constants";

const logger = setupLogger("app");

const initializeApp = async () => {
    const app = configureExpress({
        logger,
    });

    const server = await configureApolloServer();
    server.applyMiddleware({ app });

    app.get("/ping", (req, res) => {
        res.send("pong");
    });

    mockServer(app);
    app.listen(PORT, () => logger.info(`Server listening on port: ${PORT}`));
};

initializeApp().catch(error => logger.error(error));
