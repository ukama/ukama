import "reflect-metadata";
import setupLogger from "./config/logger";
import configureExpress from "./config/express";
import configureApolloServer from "./config/apolloServer";
import { mockServer } from "./mockServer";
import { PORT } from "./constants";
import { createServer } from "http";
import { job } from "./jobs/subscriptionJob";
import {
    renderPlaygroundPage,
    RenderPageOptions as PlaygroundRenderPageOptions,
} from "@apollographql/graphql-playground-html";

const logger = setupLogger("app");

const initializeApp = async () => {
    const app = configureExpress({
        logger,
    });

    const corsOption = {
        origin: ["https://*.dev.ukama.com"],
        credentials: true,
    };
    const { server, schema } = await configureApolloServer();
    server.applyMiddleware({ app, cors: corsOption });

    const httpServer = createServer(app);
    server.installSubscriptionHandlers(httpServer);

    app.get("/playground", (req, res) => {
        if (hasSession(req.headers.cookie || "")) {
            const playgroundRenderPageOptions: PlaygroundRenderPageOptions = {
                endpoint: "/graphql",
            };

            const playground = renderPlaygroundPage(
                playgroundRenderPageOptions
            );
            res.write(playground);
            res.end();
        } else {
            res.cookie("redirect", "https://bff.dev.ukama.com/playground");
            res.redirect("https://auth.dev.ukama.com/auth/login");
            res.end();
        }
    });

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

const hasSession = (session: string) => {
    if (session) {
        return session.includes("ukama_session");
    }
    return false;
};
