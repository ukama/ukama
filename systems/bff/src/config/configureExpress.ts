import express, { Express } from "express";
import winston from "winston";
import expressWinston from "express-winston";

interface Props {
    logger: winston.Logger;
}
const configureExpress = ({ logger }: Props): Express => {
    const app = express();
    app.use(express.json()); // support json encoded bodies
    app.use(express.urlencoded({ extended: true })); // support encoded bodies
    app.use(expressWinston.logger({ winstonInstance: logger }));
    return app;
};

export default configureExpress;
