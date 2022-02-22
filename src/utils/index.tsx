import { format } from "date-fns";
import { Alert_Type } from "../generated";
import { TObject } from "../types";

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

const getStatusByType = (status: string) => {
    if (status === "PENDING" || status === "BEING_CONFIGURED")
        return "Your network is being configured";
    else if (status === "ONBOARDED" || status === "ONLINE")
        return "Your network is online and well for ";
    else return "Something went wrong.";
};

const parseObjectInNameValue = (obj: any) => {
    let updatedObj: TObject[] = [];
    if (obj) {
        updatedObj = Object.keys(obj).map(key => {
            return {
                name: key,
                value:
                    key === "timestamp"
                        ? format(obj[key], "MMM dd HH:mm:ss")
                        : obj[key],
            };
        });

        let removeIndex = updatedObj
            .map(item => item?.name)
            .indexOf("__typename");
        ~removeIndex && updatedObj.splice(removeIndex, 1);
        removeIndex = updatedObj.map(item => item?.name).indexOf("id");
        ~removeIndex && updatedObj.splice(removeIndex, 1);
    }

    return updatedObj;
};

const uniqueObjectsArray = (name: string, list: TObject[]): TObject[] | [] => {
    const last =
        list.length > 0
            ? list.filter((item: TObject) => item.name !== name)
            : [];
    return last;
};

const hexToRGB = (hex: string, alpha: number): string => {
    var h = "0123456789ABCDEF";
    var r = h.indexOf(hex[1]) * 16 + h.indexOf(hex[2]);
    var g = h.indexOf(hex[3]) * 16 + h.indexOf(hex[4]);
    var b = h.indexOf(hex[5]) * 16 + h.indexOf(hex[6]);
    if (alpha) {
        return `rgba(${r}, ${g}, ${b}, ${alpha})`;
    }

    return `rgba(${r}, ${g}, ${b})`;
};

export {
    hexToRGB,
    getColorByType,
    getStatusByType,
    getTitleFromPath,
    uniqueObjectsArray,
    parseObjectInNameValue,
};
