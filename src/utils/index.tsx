import { format } from "date-fns";
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

const parseObjectInNameValue = (obj: any) => {
    obj?.id && delete obj.id;

    let updatedObj = Object.keys(obj).map(key => {
        return {
            name: key,
            value:
                key === "timestamp"
                    ? format(obj[key], "MMM dd HH:mm:ss")
                    : obj[key],
        };
    });

    const removeIndex = updatedObj.map(item => item.name).indexOf("__typename");
    ~removeIndex && updatedObj.splice(removeIndex, 1);

    return updatedObj;
};

export { getTitleFromPath, getColorByType, parseObjectInNameValue };
