import { Express } from "express";
import {
    getUser,
    getDataUsage,
    getDataBill,
    getAlerts,
    getNodes,
    getEsims,
    activateUser,
    getUsers,
    addNode,
} from "./utils";

export const mockServer = (app: Express): void => {
    app.get("/user/get_conneted_users", getUser);
    app.get("/data/data_usage", getDataUsage);
    app.get("/data/data_bill", getDataBill);
    app.get("/alert/get_alerts", getAlerts);
    app.get("/node/get_nodes", getNodes);
    app.get("/esims/get_esims", getEsims);
    app.post("/user/active_user", activateUser);
    app.get("/user/get_users", getUsers);
    app.post("/node/add_node", addNode);
};
