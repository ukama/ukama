import { Express } from "express";
import {
    getUser,
    getDataUsage,
    getDataBill,
    getAlerts,
    getNodes,
    getResidents,
} from "./utils";

export const mockServer = (app: Express): void => {
    app.get("/user/get_conneted_users", getUser);
    app.get("/data/data_usage", getDataUsage);
    app.get("/data/data_bill", getDataBill);
    app.get("/alert/get_alerts", getAlerts);
    app.get("/node/get_nodes", getNodes);
    app.get("/resident/get_residents", getResidents);
};
