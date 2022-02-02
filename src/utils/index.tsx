import { Alert_Type } from "../generated";

const getTitleFromPath = (path: string) => {
    switch (path) {
        case "/":
            return "Home";
        case "/settings":
            return "Settings";
        case "/notification":
            return "Notification";
        case "/nodes":
            return "Nodes";
        case "/users":
            return "Users";
        case "/billing":
            return "Billing";
        case "/store":
            return "Module Store";
        default:
            return "Home";
    }
};

const getColorByType = (type: Alert_Type) => {
    switch (type) {
        case Alert_Type.Error:
            return "error";
        case Alert_Type.Warning:
            return "warning";
        default:
            return "success";
    }
};

export { getTitleFromPath, getColorByType };
