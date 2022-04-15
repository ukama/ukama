import { Express } from "express";
import {
    getDataUsage,
    getDataBill,
    getAlerts,
    getNodes,
    getEsims,
    activateUser,
    getUsers,
    addNode,
    getCurrentBill,
    getBillHistory,
    getNetwork,
    updateUser,
    deleteRes,
    getUserByID,
    getNodeDetails,
    getNodeNetwork,
    getSoftwareLogs,
    getNodeApps,
} from "./utils";

export const mockServer = (app: Express): void => {
    app.get("/data/data_usage", getDataUsage);
    app.get("/data/data_bill", getDataBill);
    app.get("/alert/get_alerts", getAlerts);
    app.get("/node/get_nodes", getNodes);
    app.get("/esims/get_esims", getEsims);
    app.post("/user/active_user", activateUser);
    app.get("/user/get_users", getUsers);
    app.post("/node/add_node", addNode);
    app.get("/bill/get_current_bill", getCurrentBill);
    app.get("/bill/get_bill_history", getBillHistory);
    app.get("/network/get_network", getNetwork);
    app.post("/user/update_user", updateUser);
    app.post("/user/deactivate_user", deleteRes);
    app.post("/node/delete_node", deleteRes);
    app.get("/user/get_user", getUserByID);
    app.get("/node/node_details", getNodeDetails);
    app.get("/node/get_network", getNodeNetwork);
    app.get("/software_logs", getSoftwareLogs);
    app.get("/node_apps", getNodeApps);
};
