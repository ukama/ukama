import { Express } from "express";
import { getUser, getDataUsage } from "./utils";

export const mockServer = (app: Express): void => {
    app.get("/user/get_conneted_users", getUser);
    app.get("/data/data_usage", getDataUsage);
};
