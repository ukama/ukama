import "reflect-metadata";
import setupLogger from "./config/setupLogger";
import configureExpress from "./config/configureExpress";
import configureApolloServer from "./config/configureApolloServer";
import { mockServer } from "./mockServer";
import { PORT } from "./constants";
import { createServer } from "http";
import { job } from "./jobs/subscriptionJob";
import {
    renderPlaygroundPage,
    RenderPageOptions as PlaygroundRenderPageOptions,
} from "@apollographql/graphql-playground-html";
import cors from "cors";
import cookieParser from "cookie-parser";

const logger = setupLogger("app");

const initializeApp = async () => {
    const app = configureExpress({
        logger,
    });

    const corsOptions = {
        origin: ["https://app.dev.ukama.com", "http://localhost:3000"],
        credentials: true,
    };
    app.use(cors(corsOptions));
    app.use(cookieParser());

    const { server, schema } = await configureApolloServer();
    server.applyMiddleware({ app, cors: false });

    const httpServer = createServer(app);
    server.installSubscriptionHandlers(httpServer);

    app.get("/playground", (req, res) => {
        const playgroundRenderPageOptions: PlaygroundRenderPageOptions = {
            endpoint: "/graphql",
        };

        const playground = renderPlaygroundPage(playgroundRenderPageOptions);
        res.write(playground);
        res.end();
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
