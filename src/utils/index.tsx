import { Alert_Type } from "../generated";
import { colors } from "../theme";

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
        case "/user":
            return "User";
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
            return colors.red;
        case Alert_Type.Warning:
            return colors.yellow;
        default:
            return colors.green;
    }
};

export { getTitleFromPath, getColorByType };
