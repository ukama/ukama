import "reflect-metadata";
import setupLogger from "./config/logger";
import configureExpress from "./config/express";
import configureApolloServer from "./config/apolloServer";
import { mockServer } from "./mockServer";
import { PORT } from "./constants";
import { createServer } from "http";
import { job } from "./jobs/subscriptionJob";

const logger = setupLogger("app");

const initializeApp = async () => {
    const app = configureExpress({
        logger,
    });

    // const corsOption = { origin: ["http://localhost:3000"], credentials: true };
    const { server, schema } = await configureApolloServer();
    server.applyMiddleware({ app });

    const httpServer = createServer(app);
    server.installSubscriptionHandlers(httpServer);

    app.get("/ping", (req, res) => {
        res.send("pong");
    });

    mockServer(app);
    httpServer.listen(PORT, () => {
        logger.info(`Server listening on port: ${PORT}`);
        job(schema);
    });
};

initializeApp().catch(error => logger.error(error));
