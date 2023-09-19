import express from "express";
import expressWinston from "express-winston";

function configureExpress(logger: any) {
  const app = express();
  app.use(expressWinston.logger({ winstonInstance: logger }));
  return app;
}

export { configureExpress };
